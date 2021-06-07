package retry

import (
	"time"

	"github.com/avast/retry-go"
)

type Waiter struct {
	userCfg *WaitConfig
	libCfg  *retry.Config
	waitFn  retry.DelayTypeFunc
}

func NewWaiter(cfg *WaitConfig) *Waiter {
	waitFn := getDelayFunc(cfg.WaitType)

	opts := []retry.Option{
		retry.MaxDelay(cfg.MaxWait),
		retry.Delay(cfg.BaseWait),
		retry.MaxJitter(cfg.MaxJitter),
		retry.DelayType(waitFn),
	}

	libCfg := &retry.Config{}
	for _, opt := range opts {
		opt(libCfg)
	}

	return &Waiter{
		userCfg: cfg,
		libCfg:  libCfg,
		waitFn:  waitFn,
	}
}

// Get returns the wait duration for given iteration.
func (w *Waiter) Get(iter uint) time.Duration {
	wait := w.waitFn(iter, nil, w.libCfg)
	if wait > w.userCfg.MaxWait {
		wait = w.userCfg.MaxWait
	}
	return wait
}
