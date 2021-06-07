package promlib

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type observerVecRegister struct {
	mu    *sync.RWMutex
	known map[string]prometheus.ObserverVec
}

func (reg *observerVecRegister) get(name string) (prometheus.ObserverVec, bool) {
	reg.mu.RLock()
	sum, ok := reg.known[name]
	reg.mu.RUnlock()

	return sum, ok
}

func (reg *observerVecRegister) add(name string, sum prometheus.ObserverVec) prometheus.ObserverVec {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	exist, ok := reg.known[name]
	if ok {
		return exist
	}

	prometheus.MustRegister(sum)
	reg.known[name] = sum

	return sum
}

type observerRegister struct {
	mu    *sync.RWMutex
	known map[string]prometheus.Observer
}

func (reg *observerRegister) get(name string) (prometheus.Observer, bool) {
	reg.mu.RLock()
	hist, ok := reg.known[name]
	reg.mu.RUnlock()

	return hist, ok
}

func (reg *observerRegister) addSummary(name string, obs prometheus.Summary) prometheus.Observer {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	exist, ok := reg.known[name]
	if ok {
		return exist
	}

	prometheus.MustRegister(obs)
	reg.known[name] = obs

	return obs
}

func (reg *observerRegister) addHistogram(name string, obs prometheus.Histogram) prometheus.Observer {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	exist, ok := reg.known[name]
	if ok {
		return exist
	}

	prometheus.MustRegister(obs)
	reg.known[name] = obs

	return obs
}
