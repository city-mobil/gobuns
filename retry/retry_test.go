package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/avast/retry-go"
	"github.com/stretchr/testify/assert"
)

func TestRetrier_UnrecoverableError(t *testing.T) {
	cfg := NewDefRetryConfig()
	retrier := New(cfg)

	action := func() error {
		return retry.Unrecoverable(errors.New("fatal error"))
	}

	calls := 0
	onRetry := func(n uint, err error) {
		calls++
	}
	start := time.Now()
	err := retrier.Do(context.Background(), action, onRetry)
	assert.Error(t, err)
	assert.WithinDuration(t, time.Now(), start, cfg.BaseWait)
	assert.Equal(t, 0, calls)
}

func TestRetrier_AllAttempts(t *testing.T) {
	cfg := NewDefRetryConfig()
	retrier := New(cfg)

	action := func() error {
		return errors.New("fatal error")
	}

	calls := 0
	onRetry := func(n uint, err error) {
		calls++
	}
	start := time.Now()
	err := retrier.Do(context.Background(), action, onRetry)
	assert.WithinDuration(t, time.Now(), start, time.Duration(cfg.MaxAttempts)*cfg.BaseWait)
	assert.Error(t, err)
	assert.Equal(t, cfg.MaxAttempts, calls)
}

func TestRetrier_NoAttempts(t *testing.T) {
	cfg := NewDefRetryConfig()
	cfg.MaxAttempts = 0
	retrier := New(cfg)

	action := func() error {
		panic("this code should not be reached")
	}

	onRetry := func(n uint, err error) {}

	err := retrier.Do(context.Background(), action, onRetry)
	assert.ErrorIs(t, err, ErrNoAttempts)
}
