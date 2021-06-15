package kafka

import "time"

type statsUpdater struct {
	refreshInterval time.Duration
	stop            chan struct{}
	enabled         bool
}

func (s *statsUpdater) run(cb func()) {
	if !s.enabled {
		return
	}

	go func() {
		t := time.NewTicker(s.refreshInterval)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				cb()
			case <-s.stop:
				return
			}
		}
	}()
}
