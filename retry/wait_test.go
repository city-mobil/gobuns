package retry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBackOffWaiter(t *testing.T) {
	waiter := NewWaiter(WaitConfigWithDefaults(&WaitConfig{
		BaseWait:  DefaultBaseWait,
		MaxJitter: DefaultMaxJitter,
		MaxWait:   1 * time.Second,
		WaitType:  BackOff,
	}, time.Second))

	assert.Equal(t, DefaultBaseWait, waiter.Get(0))
	assert.Equal(t, 2*DefaultBaseWait, waiter.Get(1))
	assert.Equal(t, 1*time.Second, waiter.Get(10)) // backoff get more then maxWait
}

func TestNewFixedWaiter(t *testing.T) {
	waiter := NewWaiter(WaitConfigWithDefaults(&WaitConfig{
		BaseWait: 10 * time.Second,
		MaxWait:  1 * time.Second,
	}, time.Second))

	assert.Equal(t, 1*time.Second, waiter.Get(0)) // Fixed get more then maxWait
	assert.Equal(t, 1*time.Second, waiter.Get(1))
	assert.Equal(t, 1*time.Second, waiter.Get(10))
}

func TestNewDefaultWaiter(t *testing.T) {
	waiter := NewWaiter(NewDefWaitConfig())

	assert.Equal(t, DefaultBaseWait, waiter.Get(0))
	assert.Equal(t, DefaultBaseWait, waiter.Get(1))
	assert.Equal(t, DefaultBaseWait, waiter.Get(10))
}
