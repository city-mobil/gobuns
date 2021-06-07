package barber

import (
	"github.com/city-mobil/gobuns/config"
)

const (
	defaultThreshold = 42
	defaultMaxFails  = 50
)

type Config struct {
	// Threshold is a period of time after which
	// last error is erased.
	//
	// That means that if you have an error on the first second, then
	// this error is forgotten on first + threshold second.
	Threshold uint32

	// MaxFails indicates an allowed amount of fails to believe that host is alive.
	//
	// Hosts with total amount of errors more than MaxFails in a time of Threshold
	// will be marked as unavailable.
	MaxFails uint32
}

func NewConfig(prefix string) func() *Config {
	// TODO(a.petrukhin): implement
	prefix += "."
	var (
		threshold = config.Uint32(prefix+"threshold", defaultThreshold, "Circuit breaker closing threshold")
		maxFails  = config.Uint32(prefix+"max_fails", defaultMaxFails, "Circuit breaker max fails amount")
	)
	return func() *Config {
		return &Config{
			Threshold: *threshold,
			MaxFails:  *maxFails,
		}
	}
}

func (c *Config) withDefaults() *Config {
	if c == nil {
		c = &Config{
			Threshold: defaultThreshold,
			MaxFails:  defaultMaxFails,
		}
	}

	if c.Threshold == 0 {
		c.Threshold = defaultThreshold
	}

	if c.MaxFails == 0 {
		c.MaxFails = defaultMaxFails
	}
	return c
}
