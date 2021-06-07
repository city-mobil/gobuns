package external

import (
	"crypto/tls"
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

	cfgFn := NewConfig("http.google")

	err = config.InitOnce()
	require.NoError(t, err)

	cfg := cfgFn()

	assert.Equal(t, "ask_google", cfg.Name)
	assert.True(t, cfg.Metrics.Collect)
	assert.NotNil(t, cfg.DialContext)
	assert.Equal(t, 500, cfg.MaxIdleConns)
	assert.Equal(t, 1*time.Minute, cfg.IdleConnTimeout)
	assert.Equal(t, 200*time.Millisecond, cfg.RequestTimeout)
	assert.Equal(t, 25*time.Millisecond, cfg.RetryConfig.BaseWait)
	assert.Equal(t, 3, cfg.RetryConfig.MaxAttempts)
	assert.Equal(t, retry.BackOff, cfg.RetryConfig.WaitType)
	assert.Nil(t, cfg.OnRetry)
	assert.Equal(t, VersionTLS11, cfg.MinVersionTLS)
	assert.Equal(t, "testdata/certs/client.crt", cfg.PublicCertPath)
	assert.Equal(t, "testdata/certs/client.key", cfg.PrivateCertPath)
	assert.False(t, cfg.NoHTTPS)
}

func TestCastTLSVersion(t *testing.T) {
	tests := []struct {
		name    string
		version VersionTLS
		want    uint16
	}{
		{
			name:    "1.0",
			version: VersionTLS10,
			want:    tls.VersionTLS10,
		},
		{
			name:    "1.1",
			version: VersionTLS11,
			want:    tls.VersionTLS11,
		},
		{
			name:    "1.2",
			version: VersionTLS12,
			want:    tls.VersionTLS12,
		},
		{
			name:    "1.3",
			version: VersionTLS13,
			want:    tls.VersionTLS13,
		},
		{
			name:    "Unknown_FallbackToDefault",
			version: "5.1",
			want:    tls.VersionTLS12,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := castTLSVersion(tt.version)
			assert.Equal(t, tt.want, got)
		})
	}
}
