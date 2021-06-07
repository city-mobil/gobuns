package retry

import (
	"time"

	"github.com/avast/retry-go"
	"github.com/city-mobil/gobuns/config"
)

type WaitStrategy string

const (
	// Fixed keeps wait duration the same through all iterations.
	Fixed WaitStrategy = "fixed"
	// Random picks a random duration up to MaxJitter.
	Random WaitStrategy = "random"
	// BackOff increases duration between consecutive retries.
	BackOff WaitStrategy = "backoff"
	// Combine combines BackOff and Random strategies.
	Combine WaitStrategy = "combine"
)

func getDelayFunc(s WaitStrategy) retry.DelayTypeFunc {
	switch s {
	case BackOff:
		return retry.BackOffDelay
	case Random:
		return retry.RandomDelay
	case Combine:
		return retry.CombineDelay(retry.BackOffDelay, retry.RandomDelay)
	default:
		return retry.FixedDelay
	}
}

const (
	DefaultWaitType  = string(Fixed)
	DefaultAttempts  = 5
	DefaultBaseWait  = 10 * time.Millisecond
	DefaultMaxWait   = 10 * time.Millisecond
	DefaultMaxJitter = 10 * time.Millisecond
)

// WaitConfig is a configuration to calc wait duration.
type WaitConfig struct {
	// BaseWait is a base duration for Fixed, BackOff strategies.
	BaseWait time.Duration

	// MaxJitter sets the maximum random jitter for Random strategy.
	MaxJitter time.Duration

	// MaxWait is a maximum possible duration for any strategy.
	MaxWait time.Duration

	// WaitType is a wait strategy.
	//
	// By default: Fixed.
	// Other options:
	//   - "backoff": BackOff increases delay between consecutive retries,
	//   - "random": Random picks a random delay up to MaxJitter,
	//   - "combine": Use BackOff and Random together.
	WaitType WaitStrategy
}

// Config is a retry configuration.
type Config struct {
	WaitConfig

	// MaxAttempts is a maximum attempts to make an action.
	MaxAttempts int
}

func GetRetryConfig(prefix string) func() *Config {
	prefix = config.SanitizePrefix(prefix)
	var (
		maxAttempts = config.Int(prefix+"max_attempts", DefaultAttempts, "max action attempts")
		baseWait    = config.Duration(prefix+"base_wait", DefaultBaseWait, "base wait duration for Fixed or BackOff strategies")
		maxWait     = config.Duration(prefix+"max_wait", DefaultMaxWait, "maximum wait duration between retries")
		maxJitter   = config.Duration(prefix+"max_jitter", DefaultMaxJitter, "maximum random jitter for Random strategy")
		waitType    = config.String(prefix+"wait_type", DefaultWaitType, "wait strategy: backoff, random, combine, fixed (default)")
	)

	return func() *Config {
		conf := &Config{
			WaitConfig: WaitConfig{
				BaseWait:  *baseWait,
				MaxJitter: *maxJitter,
				MaxWait:   *maxWait,
				WaitType:  WaitStrategy(*waitType),
			},
			MaxAttempts: *maxAttempts,
		}
		return ConfigWithDefaults(conf, DefaultBaseWait)
	}
}

func GetWaitConfig(prefix string, defBaseWait time.Duration) func() *WaitConfig {
	prefix = config.SanitizePrefix(prefix)
	var (
		baseWait  = config.Duration(prefix+"base_wait", defBaseWait, "base wait duration for Fixed or BackOff strategies")
		maxWait   = config.Duration(prefix+"max_wait", DefaultMaxWait, "maximum wait duration between retries")
		maxJitter = config.Duration(prefix+"max_jitter", DefaultMaxJitter, "maximum random jitter for Random strategy")
		waitType  = config.String(prefix+"wait_type", DefaultWaitType, "wait strategy: backoff, random, combine, fixed (default)")
	)

	return func() *WaitConfig {
		conf := &WaitConfig{
			BaseWait:  *baseWait,
			MaxJitter: *maxJitter,
			MaxWait:   *maxWait,
			WaitType:  WaitStrategy(*waitType),
		}
		return WaitConfigWithDefaults(conf, defBaseWait)
	}
}

func NewDefRetryConfig() *Config {
	return ConfigWithDefaults(nil, DefaultBaseWait)
}

func NewDefWaitConfig() *WaitConfig {
	return WaitConfigWithDefaults(nil, DefaultBaseWait)
}

func ConfigWithDefaults(cfg *Config, defBaseWait time.Duration) *Config {
	if cfg == nil {
		waitCfg := WaitConfigWithDefaults(nil, defBaseWait)
		cfg = &Config{
			WaitConfig:  *waitCfg,
			MaxAttempts: DefaultAttempts,
		}
		return cfg
	}

	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = DefaultAttempts
	}

	waitCfg := WaitConfigWithDefaults(&cfg.WaitConfig, defBaseWait)
	cfg.WaitConfig = *waitCfg

	return cfg
}

func WaitConfigWithDefaults(cfg *WaitConfig, defBaseWait time.Duration) *WaitConfig {
	if cfg == nil {
		cfg = &WaitConfig{
			BaseWait:  defBaseWait,
			MaxJitter: DefaultMaxJitter,
			MaxWait:   DefaultMaxWait,
			WaitType:  WaitStrategy(DefaultWaitType),
		}
		return cfg
	}

	if cfg.WaitType == "" {
		cfg.WaitType = WaitStrategy(DefaultWaitType)
	}

	if cfg.BaseWait < 0 {
		cfg.BaseWait = defBaseWait
	}

	if cfg.MaxWait <= 0 {
		if cfg.BaseWait > DefaultMaxWait {
			cfg.MaxWait = cfg.BaseWait
		} else {
			cfg.MaxWait = DefaultMaxWait
		}
	}

	if cfg.MaxJitter <= 0 {
		cfg.MaxJitter = cfg.MaxWait
	}

	return cfg
}
