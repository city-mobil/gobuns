package registry

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

func TestNewConfig(t *testing.T) {
	cfgPath, err := filepath.Abs("testdata/config.yml")
	require.NoError(t, err)
	os.Args = append(os.Args, "--config="+cfgPath)

	cfgFn := NewConfig("app")

	err = config.InitOnce()
	require.NoError(t, err)

	cfg := cfgFn()

	assert.Equal(t, "localhost:8500", cfg.Addr)
	assert.Equal(t, 2*time.Second, cfg.QueryTimeout)
	assert.Equal(t, 25*time.Millisecond, cfg.RetryConfig.BaseWait)
	assert.Equal(t, 3, cfg.RetryConfig.MaxAttempts)
	assert.Equal(t, retry.BackOff, cfg.RetryConfig.WaitType)
}
