package mysql

import (
	"context"
	"errors"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock_mysql "github.com/city-mobil/gobuns/mocks/mysql"
	"github.com/city-mobil/gobuns/mysql/mysqlconfig"
)

var (
	tCtx   = context.Background()
	tAddr  = "127.0.0.1:3306"
	tQuery = "SELECT 1 FROM temp;"
)

type onErrCb struct {
	lastHost string
	lastErr  error
	execN    int
}

func (fn *onErrCb) run(host string, err error) {
	fn.lastHost = host
	fn.lastErr = err
	fn.execN++
}

func TestExecMasterTolerant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conn := mock_mysql.NewMockAdapter(ctrl)
	conn.EXPECT().ExecContext(tCtx, tQuery, gomock.Any()).Times(1).Return(nil, nil)

	sh, onErr := newMockedShard(conn)
	_, err := sh.ExecMasterTolerant(tCtx, tQuery)
	assert.NoError(t, err)
	assert.Empty(t, onErr.lastHost)
	assert.Nil(t, onErr.lastErr)
	assert.Empty(t, onErr.execN)
}

func TestExecMasterTolerant_PassedOnRetry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wantErr := newRetryableErr()

	conn := mock_mysql.NewMockAdapter(ctrl)
	conn.EXPECT().ExecContext(tCtx, tQuery, gomock.Any()).Times(1).Return(nil, wantErr)
	conn.EXPECT().ComponentID().Times(1).Return(tAddr)
	conn.EXPECT().ExecContext(tCtx, tQuery, gomock.Any()).Times(1).Return(nil, nil)

	sh, onErr := newMockedShard(conn)
	_, err := sh.ExecMasterTolerant(tCtx, tQuery)
	assert.NoError(t, err)
	assert.Equal(t, tAddr, onErr.lastHost)
	assert.Equal(t, wantErr, onErr.lastErr)
	assert.Equal(t, 1, onErr.execN)
}

func TestExecMasterTolerant_FailedAllRetries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wantErr := newRetryableErr()

	conn := mock_mysql.NewMockAdapter(ctrl)
	conn.EXPECT().ExecContext(tCtx, tQuery, gomock.Any()).Times(3).Return(nil, wantErr)
	conn.EXPECT().ComponentID().Times(2).Return(tAddr)

	sh, onErr := newMockedShard(conn)
	_, err := sh.ExecMasterTolerant(tCtx, tQuery)
	assert.Equal(t, wantErr, err)
	assert.Equal(t, tAddr, onErr.lastHost)
	assert.Equal(t, wantErr, onErr.lastErr)
	assert.Equal(t, 2, onErr.execN)
}

func TestExecMasterTolerant_NotRetryableErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wantErr := errors.New("fatal error")

	conn := mock_mysql.NewMockAdapter(ctrl)
	conn.EXPECT().ExecContext(tCtx, tQuery, gomock.Any()).Times(1).Return(nil, wantErr)

	sh, onErr := newMockedShard(conn)
	_, err := sh.ExecMasterTolerant(tCtx, tQuery)
	assert.Equal(t, wantErr, err)
	assert.Empty(t, onErr.lastHost)
	assert.Nil(t, onErr.lastErr)
	assert.Empty(t, onErr.execN)
}

func newMockedShard(adapter Adapter) (Shard, *onErrCb) {
	onErr := &onErrCb{}
	dbCfg := mysqlconfig.NewDefaultDatabaseConfig()
	retryCfg := &mysqlconfig.RetryConfig{
		Max:       2,
		ExecOnErr: onErr.run,
	}
	shardCfg := mysqlconfig.NewWithSingleNode(dbCfg, retryCfg)
	shard := &shard{
		config: shardCfg,
		master: adapter,
		slaves: []Adapter{adapter},
	}

	return shard, onErr
}

func newRetryableErr() error {
	return &mysql.MySQLError{
		Number:  1317,
		Message: "mock error: 1317",
	}
}
