package retry

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/avast/retry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/city-mobil/gobuns/config"
)

var (
	fullRetryCfgFn, partRetryCfgFn func() *Config
	fullWaitCfgFn, partWaitCfgFn   func() *WaitConfig
)

func init() {
	fullRetryCfgFn = GetRetryConfig("retries.full")
	partRetryCfgFn = GetRetryConfig("retries.part")
	fullWaitCfgFn = GetWaitConfig("wait.full", DefaultBaseWait)
	partWaitCfgFn = GetWaitConfig("wait.part", DefaultBaseWait)
}

func TestNewFullRetryConfig(t *testing.T) {
	configPath, err := filepath.Abs("testdata/retry.yml")
	require.NoError(t, err)

	os.Args = append(os.Args, "--config="+configPath)
	err = config.InitOnce()
	require.NoError(t, err)

	cfg := fullRetryCfgFn()
	assert.NotNil(t, cfg)

	assert.Equal(t, 3, cfg.MaxAttempts)
	assert.Equal(t, Random, cfg.WaitType)
	assert.Equal(t, 2*time.Second, cfg.MaxJitter)
	assert.Equal(t, DefaultBaseWait, cfg.BaseWait)
	assert.Equal(t, 1*time.Second, cfg.MaxWait)

	retrier := New(cfg)

	action := func() error {
		return retry.Unrecoverable(errors.New("fatal error"))
	}

	onRetry := func(n uint, err error) {
	}
	start := time.Now()
	err = retrier.Do(context.Background(), action, onRetry)
	require.Error(t, err)
	require.WithinDuration(t, time.Now(), start, 1*time.Second)
}

func TestNewPartRetryConfig(t *testing.T) {
	configPath, err := filepath.Abs("testdata/retry.yml")
	require.NoError(t, err)

	os.Args = append(os.Args, "--config="+configPath)
	err = config.InitOnce()
	require.NoError(t, err)

	cfg := partRetryCfgFn()
	assert.NotNil(t, cfg)

	assert.Equal(t, DefaultAttempts, cfg.MaxAttempts)
	assert.Equal(t, Fixed, cfg.WaitType)
	assert.Equal(t, DefaultMaxJitter, cfg.MaxJitter)
	assert.Equal(t, 10*time.Millisecond, cfg.BaseWait)
	assert.Equal(t, 1*time.Second, cfg.MaxWait)

	retrier := New(cfg)

	action := func() error {
		return errors.New("fatal error")
	}

	onRetry := func(n uint, err error) {
	}
	start := time.Now().UnixNano()
	err = retrier.Do(context.Background(), action, onRetry)
	require.Error(t, err)
	require.True(t, time.Now().UnixNano()-start >= (DefaultAttempts-1)*10*time.Millisecond.Nanoseconds())
}

func TestNewFullWaitConfig(t *testing.T) {
	configPath, err := filepath.Abs("testdata/retry.yml")
	require.NoError(t, err)

	os.Args = append(os.Args, "--config="+configPath)
	err = config.InitOnce()
	require.NoError(t, err)

	cfg := fullWaitCfgFn()
	assert.NotNil(t, cfg)

	assert.Equal(t, BackOff, cfg.WaitType)
	assert.Equal(t, 2*time.Second, cfg.MaxJitter)
	assert.Equal(t, DefaultBaseWait, cfg.BaseWait)
	assert.Equal(t, 1*time.Second, cfg.MaxWait)

	waiter := NewWaiter(cfg)

	assert.Equal(t, DefaultBaseWait, waiter.Get(0))
	assert.Equal(t, 2*DefaultBaseWait, waiter.Get(1))
	assert.Equal(t, 1*time.Second, waiter.Get(10)) // backoff get more then maxWait
}

func TestNewPartWaitConfig(t *testing.T) {
	configPath, err := filepath.Abs("testdata/retry.yml")
	require.NoError(t, err)

	os.Args = append(os.Args, "--config="+configPath)
	err = config.InitOnce()
	require.NoError(t, err)

	cfg := partWaitCfgFn()
	assert.NotNil(t, cfg)

	assert.Equal(t, Fixed, cfg.WaitType)
	assert.Equal(t, DefaultMaxJitter, cfg.MaxJitter)
	assert.Equal(t, 10*time.Second, cfg.BaseWait)
	assert.Equal(t, 1*time.Second, cfg.MaxWait)

	waiter := NewWaiter(cfg)

	assert.Equal(t, 1*time.Second, waiter.Get(0)) // Fixed get more then maxWait
	assert.Equal(t, 1*time.Second, waiter.Get(1))
	assert.Equal(t, 1*time.Second, waiter.Get(10))
}
