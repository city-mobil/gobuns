package promlib

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	eventMutex = &sync.RWMutex{}
)

var (
	registeredCounter    = make(map[string]prometheus.Counter)
	registeredCounterVec = make(map[string]*prometheus.CounterVec)
)

type Event struct {
	Name        string
	Namespace   string
	Subsystem   string
	Help        string
	ConstLabels Labels
}

type Labels = map[string]string

// IncCnt creates a counter in place using the given name
// and global application name, increments the counter by 1.
func IncCnt(name string) {
	e := newEventInPlace(name)
	IncCntEvent(e)
}

func IncCntWithLabels(name string, params Labels) {
	e := newEventInPlace(name)
	IncCntEventWithLabels(e, params)
}

// AddCnt creates a counter in place using the given name
// and global application name, adds the value to the counter.
func AddCnt(name string, val float64) {
	e := newEventInPlace(name)
	AddCntEvent(e, val)
}

func AddCntWithLabels(name string, params Labels, val float64) {
	e := newEventInPlace(name)
	AddCntEventWithLabels(e, params, val)
}

func IncCntEvent(e *Event) {
	counter := registerCounter(e)
	counter.Inc()
}

func AddCntEvent(e *Event, val float64) {
	counter := registerCounter(e)
	counter.Add(val)
}

func registerCounter(e *Event) prometheus.Counter {
	fqn := buildEventFQName(e)

	eventMutex.RLock()
	counter, ok := registeredCounter[fqn]
	eventMutex.RUnlock()
	if ok {
		return counter
	}

	ns := e.Namespace
	if ns == "" {
		ns = globalNS
	}
	counter = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        e.Name,
		Namespace:   ns,
		Subsystem:   e.Subsystem,
		Help:        e.Help,
		ConstLabels: e.ConstLabels,
	})

	eventMutex.Lock()
	defer eventMutex.Unlock()
	cnt, ok := registeredCounter[fqn]
	if ok {
		return cnt
	}

	prometheus.MustRegister(counter)
	registeredCounter[fqn] = counter

	return counter
}

func IncCntEventWithLabels(e *Event, params Labels) {
	counterVec := registerCounterVec(e, params)
	counterVec.With(params).Inc()
}

func AddCntEventWithLabels(e *Event, params Labels, val float64) {
	counterVec := registerCounterVec(e, params)
	counterVec.With(params).Add(val)
}

func registerCounterVec(e *Event, params Labels) *prometheus.CounterVec {
	fqn := buildEventFQName(e)

	eventMutex.RLock()
	counter, ok := registeredCounterVec[fqn]
	eventMutex.RUnlock()
	if ok {
		return counter
	}

	labels := make([]string, 0, len(params))
	for key := range params {
		labels = append(labels, key)
	}

	ns := e.Namespace
	if ns == "" {
		ns = globalNS
	}
	counter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        e.Name,
		Namespace:   ns,
		Subsystem:   e.Subsystem,
		Help:        e.Help,
		ConstLabels: e.ConstLabels,
	}, labels)

	eventMutex.Lock()
	defer eventMutex.Unlock()
	cnt, ok := registeredCounterVec[fqn]
	if ok {
		return cnt
	}

	prometheus.MustRegister(counter)
	registeredCounterVec[fqn] = counter

	return counter
}

func newEventInPlace(name string) *Event {
	return &Event{
		Name:      name,
		Namespace: globalNS,
	}
}

func buildEventFQName(e *Event) string {
	return prometheus.BuildFQName(e.Namespace, e.Subsystem, e.Name)
}
