package promlib

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type TransactionType int

const (
	Summary TransactionType = iota
	Histogram
)

var (
	DefBuckets    = prometheus.DefBuckets
	DefObjectives = map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}
)

// LinearBuckets creates 'count' buckets, each 'width' wide, where the lowest
// bucket has an upper bound of 'start'. The final +Inf bucket is not counted
// and not included in the returned slice. The returned slice is meant to be
// used for the Buckets field of HistogramOpts.
//
// The function panics if 'count' is zero or negative.
func LinearBuckets(start, width float64, count int) []float64 {
	return prometheus.LinearBuckets(start, width, count)
}

// ExponentialBuckets creates 'count' buckets, where the lowest bucket has an
// upper bound of 'start' and each following bucket's upper bound is 'factor'
// times the previous bucket's upper bound. The final +Inf bucket is not counted
// and not included in the returned slice. The returned slice is meant to be
// used for the Buckets field of HistogramOpts.
//
// The function panics if 'count' is 0 or negative, if 'start' is 0 or negative,
// or if 'factor' is less than or equal 1.
func ExponentialBuckets(start, factor float64, count int) []float64 {
	return prometheus.ExponentialBuckets(start, factor, count)
}

var (
	observers = &observerRegister{
		mu:    &sync.RWMutex{},
		known: make(map[string]prometheus.Observer),
	}

	observerVecs = &observerVecRegister{
		mu:    &sync.RWMutex{},
		known: make(map[string]prometheus.ObserverVec),
	}
)

type TransactionOpts interface {
	Type() TransactionType
	Meta() MetaOpts
	WithLabels() bool
}

type MetaOpts struct {
	// Namespace, Subsystem, and Name are components of the fully-qualified
	// name of the Transaction (created by joining these components with
	// "_"). Only Name is mandatory, the others merely help structuring the
	// name. Note that the fully-qualified name of the Transaction must be a
	// valid Prometheus metric name.
	Namespace string
	Subsystem string
	Name      string

	// Help provides information about this metric.
	//
	// Metrics with the same fully-qualified name must have the same Help
	// string.
	Help string

	// ConstLabels are used to attach fixed labels to this metric. Metrics
	// with the same fully-qualified name must have the same label names in
	// their ConstLabels.
	//
	// ConstLabels are only used rarely. In particular, do not use them to
	// attach the same labels to all your metrics. Those use cases are
	// better covered by target labels set by the scraping Prometheus
	// server, or by one specific metric (e.g. a build_info or a
	// machine_role metric). See also
	// https://prometheus.io/docs/instrumenting/writing_exporters/#target-labels-not-static-scraped-labels
	ConstLabels Labels
}

// HistogramOpts bundles the specific options for creating a Histogram metric.
//
// See prometheus docs for more information.
type HistogramOpts struct {
	MetaOpts
	Buckets []float64
	Labels  []string
}

func (o *HistogramOpts) Type() TransactionType {
	return Histogram
}

func (o *HistogramOpts) Meta() MetaOpts {
	return o.MetaOpts
}

func (o *HistogramOpts) WithLabels() bool {
	return len(o.Labels) > 0
}

func newHistogram(opts *HistogramOpts) prometheus.Histogram {
	buckets := opts.Buckets
	if buckets == nil {
		buckets = DefBuckets
	}

	return prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      opts.Name,
		Namespace: opts.Namespace,
		Subsystem: opts.Subsystem,
		Help:      opts.Help,
		Buckets:   buckets,
	})
}

func newHistogramVec(opts *HistogramOpts) *prometheus.HistogramVec {
	buckets := opts.Buckets
	if buckets == nil {
		buckets = DefBuckets
	}

	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      opts.Name,
		Namespace: opts.Namespace,
		Subsystem: opts.Subsystem,
		Help:      opts.Help,
		Buckets:   buckets,
	}, opts.Labels)
}

// SummaryOpts bundles the specific options for creating a Summary metric.
//
// See prometheus docs for more information.
type SummaryOpts struct {
	MetaOpts
	Objectives map[float64]float64
	MaxAge     time.Duration
	AgeBuckets uint32
	Labels     []string
}

func (o *SummaryOpts) Type() TransactionType {
	return Summary
}

func (o *SummaryOpts) Meta() MetaOpts {
	return o.MetaOpts
}

func (o *SummaryOpts) WithLabels() bool {
	return len(o.Labels) > 0
}

func newSummary(opts *SummaryOpts) prometheus.Summary {
	objectives := opts.Objectives
	if objectives == nil {
		objectives = DefObjectives
	}

	return prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace:  opts.Namespace,
		Subsystem:  opts.Subsystem,
		Name:       opts.Name,
		Help:       opts.Help,
		Objectives: objectives,
		MaxAge:     opts.MaxAge,
		AgeBuckets: opts.AgeBuckets,
	})
}

func newSummaryVec(opts *SummaryOpts) *prometheus.SummaryVec {
	objectives := opts.Objectives
	if objectives == nil {
		objectives = DefObjectives
	}

	return prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  opts.Namespace,
		Subsystem:  opts.Subsystem,
		Name:       opts.Name,
		Help:       opts.Help,
		Objectives: objectives,
		MaxAge:     opts.MaxAge,
		AgeBuckets: opts.AgeBuckets,
	}, opts.Labels)
}

type Transaction interface {
	Start(labels ...string)
	Observe(value float64, labels ...string)
	End()
}

func NewTransaction(opts TransactionOpts) Transaction {
	if opts == nil {
		return nil
	}

	var tx Transaction
	if opts.WithLabels() {
		tx = newTxnVec(opts)
	} else {
		tx = newTxn(opts)
	}

	return tx
}

func NewNopTransaction() Transaction {
	return &txn{}
}

type txn struct {
	observer prometheus.Observer
	timer    *prometheus.Timer
}

func newTxn(opts TransactionOpts) *txn {
	name := buildFQName(opts.Meta())

	observer, ok := observers.get(name)
	if !ok {
		observer = newObserver(name, opts)
		if observer == nil {
			return nil
		}
	}

	return &txn{
		observer: observer,
	}
}

func (t *txn) Start(_ ...string) {
	if t == nil || t.observer == nil {
		return
	}
	t.timer = prometheus.NewTimer(t.observer)
}

func (t *txn) Observe(value float64, _ ...string) {
	if t.observer != nil {
		t.observer.Observe(value)
	}
}

func (t *txn) End() {
	if t == nil || t.timer == nil {
		return
	}
	_ = t.timer.ObserveDuration()
}

type txnVec struct {
	observer prometheus.ObserverVec
	timer    *prometheus.Timer
}

func newTxnVec(opts TransactionOpts) *txnVec {
	name := buildFQName(opts.Meta())

	observer, ok := observerVecs.get(name)
	if !ok {
		observer = newObserverVec(name, opts)
		if observer == nil {
			return nil
		}
	}

	return &txnVec{
		observer: observer,
	}
}

func (t *txnVec) Start(values ...string) {
	if t == nil || t.observer == nil {
		return
	}
	t.timer = prometheus.NewTimer(t.observer.WithLabelValues(values...))
}

func (t *txnVec) Observe(value float64, labels ...string) {
	if t.observer != nil && len(labels) > 0 {
		t.observer.WithLabelValues(labels...).Observe(value)
	}
}

func (t *txnVec) End() {
	if t == nil || t.timer == nil {
		return
	}
	_ = t.timer.ObserveDuration()
}

func newObserver(name string, opts TransactionOpts) prometheus.Observer {
	switch opts.Type() {
	case Summary:
		sum := newSummary(opts.(*SummaryOpts))
		return observers.addSummary(name, sum)
	case Histogram:
		hist := newHistogram(opts.(*HistogramOpts))
		return observers.addHistogram(name, hist)
	default:
		return nil
	}
}

func newObserverVec(name string, opts TransactionOpts) prometheus.ObserverVec {
	var ob prometheus.ObserverVec

	switch opts.Type() {
	case Summary:
		ob = newSummaryVec(opts.(*SummaryOpts))
	case Histogram:
		ob = newHistogramVec(opts.(*HistogramOpts))
	default:
		return nil
	}

	return observerVecs.add(name, ob)
}

func buildFQName(opts MetaOpts) string {
	return prometheus.BuildFQName(opts.Namespace, opts.Subsystem, opts.Name)
}
