package rabbit

import (
	"bytes"
	"testing"
	"time"

	"github.com/city-mobil/gobuns/zlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnector_SuccessConnect(t *testing.T) {
	_, err := NewConnector(dummyLogger, Config{
		Addr:           "rabbitmq:5672",
		Login:          "guest",
		Password:       "guest",
		ConnectTimeout: time.Second,
	})
	require.NoError(t, err)
}

func TestConnector_ConnectErrors(t *testing.T) {
	tests := []struct {
		name             string
		addr             string
		expectedContains []string
	}{
		{
			name:             "IOTimeout",
			addr:             "2.3.4.5:5672",
			expectedContains: []string{"dial tcp 2.3.4.5:5672", "i/o timeout"},
		},
		{
			name:             "ConnectionTimeout",
			addr:             "localhost:2701",
			expectedContains: []string{"dial tcp 127.0.0.1:2701", "connect: connection refused"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConnector(dummyLogger, Config{
				Addr:           tt.addr,
				Login:          "guest",
				Password:       "guest",
				ReconnectDelay: time.Second,
				ConnectTimeout: 500 * time.Millisecond,
			})
			require.NotNil(t, err)
			for _, want := range tt.expectedContains {
				assert.Contains(t, err.Error(), want)
			}
		})
	}
}

func TestClose(t *testing.T) {
	buf := bytes.Buffer{}
	conn, err := NewConnector(zlog.New(&buf), Config{
		Addr:           "rabbitmq:5672",
		Login:          "guest",
		Password:       "guest",
		ConnectTimeout: time.Second,
		ReconnectDelay: 10 * time.Millisecond,
	})
	require.NoError(t, err)

	err = conn.Close()
	require.NoError(t, err)
	assert.Eventually(t, func() bool {
		return assert.Contains(t, buf.String(), "closing connector") &&
			assert.NotContains(t, buf.String(), "connection reconnected")
	}, time.Second, 200*time.Millisecond)
}
