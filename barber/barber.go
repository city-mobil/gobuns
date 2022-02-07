package barber

import (
	"math"
	"sync"
	"time"
)

// Barber is fast and easy to use circuit-breaker implementation
type Barber interface {
	IsAvailable(serverID int, tm time.Time) bool
	AddError(serverID int, tm time.Time)
	Stats() *Stats
}

// fail describes a single registered fail.
//
// No reason is recorded for performance reasons.
type fail struct {
	lastTS int64
	count  uint32
}

// host describes a single given host for circuit breaker.
//
// A single host is described in NewBarber initialization.
type host struct {
	mu    sync.RWMutex
	fails []*fail
}

func (h *host) countFails(ts int64, maxAllowed uint32) (res uint32) {
	timeThreshold := int64(len(h.fails))
	h.mu.RLock()
	defer h.mu.RUnlock()

	for i := 0; i < len(h.fails); i++ {
		if ts-h.fails[i].lastTS < timeThreshold {
			res += h.fails[i].count
		}
		if res > maxAllowed {
			break
		}
	}

	return
}

// addFails counts a single fail in according timestamp.
func (h *host) addFail(ts int64) {
	idx := ts % int64(len(h.fails))
	h.mu.Lock()
	defer h.mu.Unlock()

	h1 := h.fails[idx]
	if h1.lastTS != ts {
		h1.lastTS = ts
		h1.count = 1
	} else {
		h1.count++
	}
}

type barber struct {
	config *Config
	hosts  map[int]*host
}

// NewBarber creates new cirulnik-barber for a given hosts list.
func NewBarber(hosts []int, config *Config) Barber {
	config = config.withDefaults()

	mp := make(map[int]*host)
	for _, v := range hosts {
		mp[v] = &host{
			fails: make([]*fail, config.Threshold),
		}
		for j := 0; j < len(mp[v].fails); j++ {
			mp[v].fails[j] = &fail{}
		}
	}

	return &barber{
		config: config,
		hosts:  mp,
	}
}

// IsAvailable returns the availability status of host
func (b *barber) IsAvailable(serverID int, tm time.Time) bool {
	h, ok := b.hosts[serverID]
	if !ok {
		return false
	}

	res := h.countFails(tm.Unix(), b.config.MaxFails)
	return res <= b.config.MaxFails
}

// AddError adds error to the selected host with given timestamp
func (b *barber) AddError(serverID int, tm time.Time) {
	h, ok := b.hosts[serverID]
	if !ok {
		return
	}

	h.addFail(tm.Unix())
}

// Stats returns error statistics for all hosts
func (b *barber) Stats() *Stats {
	stats := &Stats{}
	stats.Hosts = make([]*StatHost, 0, len(b.hosts))
	now := time.Now().Unix()
	for k, v := range b.hosts {
		s := &StatHost{}
		s.ServerID = k
		s.FailsCount = int(v.countFails(now, math.MaxUint32))

		stats.Hosts = append(stats.Hosts, s)
	}
	return stats
}
