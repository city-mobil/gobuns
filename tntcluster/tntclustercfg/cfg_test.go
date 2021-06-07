package tntclustercfg

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/city-mobil/gobuns/config"
	"github.com/city-mobil/gobuns/retry"
)

func TestNewClusterConfig(t *testing.T) {
	configPath, err := filepath.Abs("testdata/cluster.yml")
	require.NoError(t, err)

	cfgFn := NewClusterConfig("auth")

	os.Args = append(os.Args, "--config="+configPath)
	err = config.InitOnce()
	require.NoError(t, err)

	cfg, err := cfgFn()
	require.NoError(t, err)
	assert.Len(t, cfg.Shards, 2)

	shard := cfg.Shards[0]
	assert.Equal(t, "primary", shard.Name)
	assert.Equal(t, "auth1.ddk:3301", shard.MasterAddr)
	assert.Equal(t, []string{"auth-slave1.ddk:3301", "auth-slave2.ddk:3301"}, shard.SlaveAddrs)
	assert.Equal(t, 30, shard.PoolSize)

	opts := shard.Opts
	assert.Equal(t, 1*time.Minute, opts.ConnectTimeout)
	assert.Equal(t, "clients", opts.DefaultSpace)
	assert.Equal(t, "guest", opts.User)
	assert.Equal(t, "guest", opts.Password)
	assert.Equal(t, 10, opts.PoolMaxPacketSize)

	retryCfg := shard.RetryConfig
	assert.Equal(t, 3, retryCfg.MaxAttempts)
	assert.Equal(t, 10*time.Millisecond, retryCfg.BaseWait)

	queryTimeoutCfg := shard.QueryTimeoutConfig
	assert.Equal(t, 1*time.Minute, queryTimeoutCfg.BaseWait)
	assert.Equal(t, retry.Fixed, queryTimeoutCfg.WaitType)

	shard = cfg.Shards[1]
	assert.Equal(t, "auth2.ddk:3301", shard.MasterAddr)

	opts = shard.Opts
	assert.Equal(t, "guest", opts.User)
	assert.Equal(t, "guest", opts.Password)

	retryCfg = shard.RetryConfig
	assert.Equal(t, 5, retryCfg.MaxAttempts)
	assert.Equal(t, retry.Random, retryCfg.WaitType)
	assert.Equal(t, 200*time.Millisecond, retryCfg.MaxJitter)

	queryTimeoutCfg = shard.QueryTimeoutConfig
	assert.Equal(t, retry.Random, queryTimeoutCfg.WaitType)
	assert.Equal(t, 200*time.Millisecond, queryTimeoutCfg.MaxJitter)
}
