package registry

import (
	"time"

	"github.com/city-mobil/gobuns/config"
	"github.com/city-mobil/gobuns/retry"
)

const (
	defaultAddr    = "127.0.0.1:8500"
	defaultTimeout = time.Second
)

type Config struct {
	Addr         string
	QueryTimeout time.Duration
	// RetryConfig is a configuration for retry policy.
	RetryConfig *retry.Config
}

func NewConfig(prefix string) func() Config {
	if prefix != "" {
		prefix += ".registry."
	} else {
		prefix = "registry."
	}

	var (
		addr         = config.String(prefix+"addr", defaultAddr, "consul agent address")
		queryTimeout = config.Duration(prefix+"query_timeout", defaultTimeout, "timeout for KV request")
		retryCfgFn   = retry.GetRetryConfig(prefix + "retries")
	)

	return func() Config {
		return Config{
			Addr:         *addr,
			QueryTimeout: *queryTimeout,
			RetryConfig:  retryCfgFn(),
		}
	}
}
