package zlog

import (
	"time"

	"github.com/rs/zerolog"
)

var (
	// Often samples log every ~ 10 events.
	Often = NewRandomSampler(10)
	// Sometimes samples log every ~ 100 events.
	Sometimes = NewRandomSampler(100)
	// Rarely samples log every ~ 1000 events.
	Rarely = NewRandomSampler(1000)
)

// Sampler defines an interface to a log sampler.
type Sampler interface {
	// Sample returns true if the event should be part of the sample, false if
	// the event should be dropped.
	Sample(lvl Level) bool
}

type randomSampler struct {
	z zerolog.RandomSampler
}

// NewRandomSampler returns a sampler uses a PRNG
// to randomly sample an event out of N events, regardless of their level.
func NewRandomSampler(n uint32) Sampler {
	return &randomSampler{
		z: zerolog.RandomSampler(n),
	}
}

func (s *randomSampler) Sample(lvl Level) bool {
	return s.z.Sample(lvl)
}

type basicSampler struct {
	z zerolog.BasicSampler
}

// NewBasicSampler returns a sampler that will send
// every Nth events, regardless of there level.
func NewBasicSampler(n uint32) Sampler {
	return &basicSampler{
		z: zerolog.BasicSampler{
			N: n,
		},
	}
}

func (s *basicSampler) Sample(lvl Level) bool {
	return s.z.Sample(lvl)
}

type burstSampler struct {
	z zerolog.BurstSampler
}

// NewBurstSample returns a sampler which lets Burst events pass per Period then pass the decision to
// Next sampler. If next is not set, all subsequent events are rejected.
func NewBurstSample(burst uint32, period time.Duration, next Sampler) Sampler {
	return &burstSampler{
		z: zerolog.BurstSampler{
			Burst:       burst,
			Period:      period,
			NextSampler: next,
		},
	}
}

func (s *burstSampler) Sample(lvl Level) bool {
	return s.z.Sample(lvl)
}

// LevelSampler applies a different sampler for each level.
type LevelSampler struct {
	TraceSampler, DebugSampler, InfoSampler, WarnSampler, ErrorSampler Sampler
}

func (s *LevelSampler) Sample(lvl Level) bool {
	switch lvl {
	case TraceLevel:
		if s.TraceSampler != nil {
			return s.TraceSampler.Sample(lvl)
		}
	case DebugLevel:
		if s.DebugSampler != nil {
			return s.DebugSampler.Sample(lvl)
		}
	case InfoLevel:
		if s.InfoSampler != nil {
			return s.InfoSampler.Sample(lvl)
		}
	case WarnLevel:
		if s.WarnSampler != nil {
			return s.WarnSampler.Sample(lvl)
		}
	case ErrorLevel:
		if s.ErrorSampler != nil {
			return s.ErrorSampler.Sample(lvl)
		}
	}
	return true
}
