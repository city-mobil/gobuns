package tntcluster

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/viciious/go-tarantool"

	mock_pool "github.com/city-mobil/gobuns/mocks/tntcluster/pool"
	"github.com/city-mobil/gobuns/retry"
	"github.com/city-mobil/gobuns/tntcluster/pool"
)

func TestRetryOnRetryableErrors(t *testing.T) {
	oldTntRetryableErrors := tntRetryableErrors
	defer func() {
		tntRetryableErrors = oldTntRetryableErrors
	}()
	maxAttempts := 3

	tntRetryableErrors = []uint{
		tarantool.ErrAccessDenied,
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	box, err := tarantool.NewBox("", nil)
	assert.NoError(t, err)
	defer box.Close()

	conn, err := box.Connect(&tarantool.Options{})
	assert.NoError(t, err)
	assert.NotNil(t, conn)
	defer conn.Close()

	rootSpan := opentracing.StartSpan("rootSpan")

	query := &tarantool.Call{
		Name: "unknownFunction",
	}

	tests := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "Context with OT span",
			ctx:  opentracing.ContextWithSpan(context.Background(), rootSpan),
		},
		{
			name: "Background context",
			ctx:  context.Background(),
		},
	}

	for _, tt := range tests {
		tc := tt

		mockMasterPool := mock_pool.NewMockConnectorPool(mockCtrl)
		mockMasterPool.EXPECT().Connect().Return(conn, nil).Times(maxAttempts)
		mockMasterPool.EXPECT().RemoteAddr().Return(box.Addr()).AnyTimes()

		mockSlavePool := mock_pool.NewMockConnectorPool(mockCtrl)
		mockSlavePool.EXPECT().Connect().Return(conn, nil).Times(maxAttempts)
		mockSlavePool.EXPECT().RemoteAddr().Return(box.Addr()).AnyTimes()

		retryCfg := retry.ConfigWithDefaults(&retry.Config{
			MaxAttempts: maxAttempts,
		}, 100*time.Millisecond)

		waitCfg := retry.WaitConfigWithDefaults(nil, 100*time.Millisecond)

		shard := &shard{
			options:         tarantool.Options{},
			masterConnector: mockMasterPool,
			slaveConnectors: []pool.ConnectorPool{mockSlavePool},
			waiter:          retry.NewWaiter(waitCfg),
			retrier:         retry.New(retryCfg),
			stop:            make(chan struct{}),
		}

		t.Run(tt.name, func(t *testing.T) {
			result, err := shard.CallTolerant(tc.ctx, query)
			assert.Error(t, err)
			assert.Nil(t, result)

			result, err = shard.CallReplicaTolerant(tc.ctx, query)
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}
