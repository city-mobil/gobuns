package retry

import (
	"context"
	"errors"

	"github.com/avast/retry-go"
)

var (
	ErrNoAttempts = errors.New("at least one attempt must be set to execute the action")
)

// Action is a user function executed by retry policy.
type Action = func() error

// OnRetryFunc is a function executed before every retry.
type OnRetryFunc = func(n uint, err error)

type Retrier struct {
	cfg  *Config
	opts []retry.Option

	// delayFn is called to return the next delay to wait after
	// the retriable function fails on `err` after `n` attempts.
	delayFn retry.DelayTypeFunc
}

func New(cfg *Config) *Retrier {
	delayFn := getDelayFunc(cfg.WaitType)

	opts := []retry.Option{
		retry.MaxDelay(cfg.MaxWait),
		retry.Attempts(uint(cfg.MaxAttempts)),
		retry.Delay(cfg.BaseWait),
		retry.MaxJitter(cfg.MaxJitter),
		retry.DelayType(delayFn),
	}

	return &Retrier{
		cfg:     cfg,
		opts:    opts,
		delayFn: delayFn,
	}
}

// Do executes action and retries it in case of recoverable errors.
func (p *Retrier) Do(ctx context.Context, action Action, onRetry OnRetryFunc) error {
	if p.cfg.MaxAttempts < 1 {
		return ErrNoAttempts
	}

	opts := make([]retry.Option, 0, len(p.opts)+2)
	opts = append(opts, retry.Context(ctx), retry.OnRetry(onRetry))
	opts = append(opts, p.opts...)

	return retry.Do(action, opts...)
}

func Unrecoverable(err error) error {
	return retry.Unrecoverable(err)
}
