package promlib

import (
	"fmt"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestNopTransaction(_ *testing.T) {
	nop := NewNopTransaction()
	nop.Start()
	nop.End()
}

func TestSummary(t *testing.T) {
	tests := []struct {
		name   string
		opts   *SummaryOpts
		values []string
	}{
		{
			name: "WithLabels",
			opts: &SummaryOpts{
				MetaOpts: MetaOpts{
					Namespace: "test",
					Subsystem: "unit",
					Name:      "summary_with_labels",
					Help:      "ensure that summary vec is working",
				},
				Objectives: DefObjectives,
				MaxAge:     1 * time.Minute,
				AgeBuckets: 30,
				Labels:     []string{"t1"},
			},
			values: []string{"v1"},
		},
		{
			name: "WithoutLabels",
			opts: &SummaryOpts{
				MetaOpts: MetaOpts{
					Namespace: "test",
					Subsystem: "unit",
					Name:      "summary_without_labels",
					Help:      "ensure that summary is working",
				},
				Objectives: DefObjectives,
				MaxAge:     1 * time.Minute,
				AgeBuckets: 30,
				Labels:     nil,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			txn := NewTransaction(tt.opts)
			txn.Start(tt.values...)
			time.Sleep(10 * time.Millisecond)
			txn.End()

			labels := ""
			for i := 0; i < len(tt.opts.Labels); i++ {
				labels += fmt.Sprintf(`%s="%s"`, tt.opts.Labels[i], tt.values[i])
				if i < len(tt.opts.Labels)-1 {
					labels += ","
				}
			}

			fqName := buildFQName(tt.opts.MetaOpts)

			metrics := dump(prometheus.DefaultGatherer)
			assert.Contains(t, metrics, tt.opts.MetaOpts.Help)
			if len(tt.opts.Labels) > 0 {
				assert.Contains(t, metrics, fmt.Sprintf(`%s{%s,quantile="0.5"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s{%s,quantile="0.9"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s{%s,quantile="0.99"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_count{%s}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_sum{%s}`, fqName, labels))
			} else {
				assert.Contains(t, metrics, fmt.Sprintf(`%s{quantile="0.5"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s{quantile="0.9"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s{quantile="0.99"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_count`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_sum`, fqName))
			}
		})
	}
}

func TestSummaryObserve(t *testing.T) {
	tests := []struct {
		name   string
		opts   *SummaryOpts
		values []string
	}{
		{
			name: "WithLabels",
			opts: &SummaryOpts{
				MetaOpts: MetaOpts{
					Namespace: "test",
					Subsystem: "unit",
					Name:      "summary_with_labels",
					Help:      "ensure that summary vec is working",
				},
				Objectives: DefObjectives,
				MaxAge:     1 * time.Minute,
				AgeBuckets: 30,
				Labels:     []string{"t1"},
			},
			values: []string{"v1"},
		},
		{
			name: "WithoutLabels",
			opts: &SummaryOpts{
				MetaOpts: MetaOpts{
					Namespace: "test",
					Subsystem: "unit",
					Name:      "summary_without_labels",
					Help:      "ensure that summary is working",
				},
				Objectives: DefObjectives,
				MaxAge:     1 * time.Minute,
				AgeBuckets: 30,
				Labels:     nil,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			txn := NewTransaction(tt.opts)
			txn.Observe(0.1, tt.values...)

			labels := ""
			for i := 0; i < len(tt.opts.Labels); i++ {
				labels += fmt.Sprintf(`%s="%s"`, tt.opts.Labels[i], tt.values[i])
				if i < len(tt.opts.Labels)-1 {
					labels += ","
				}
			}

			fqName := buildFQName(tt.opts.MetaOpts)

			metrics := dump(prometheus.DefaultGatherer)
			assert.Contains(t, metrics, tt.opts.MetaOpts.Help)
			if len(tt.opts.Labels) > 0 {
				assert.Contains(t, metrics, fmt.Sprintf(`%s{%s,quantile="0.5"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s{%s,quantile="0.9"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s{%s,quantile="0.99"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_count{%s}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_sum{%s}`, fqName, labels))
			} else {
				assert.Contains(t, metrics, fmt.Sprintf(`%s{quantile="0.5"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s{quantile="0.9"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s{quantile="0.99"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_count`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_sum`, fqName))
			}
		})
	}
}

func TestHistogram(t *testing.T) {
	tests := []struct {
		name   string
		opts   *HistogramOpts
		values []string
	}{
		{
			name: "WithLabels",
			opts: &HistogramOpts{
				MetaOpts: MetaOpts{
					Namespace: "test",
					Subsystem: "unit",
					Name:      "histogram_with_labels",
					Help:      "ensure that histogram vec is working",
				},
				Buckets: DefBuckets,
				Labels:  []string{"t1"},
			},
			values: []string{"v1"},
		},
		{
			name: "WithoutLabels",
			opts: &HistogramOpts{
				MetaOpts: MetaOpts{
					Namespace: "test",
					Subsystem: "unit",
					Name:      "histogram_without_labels",
					Help:      "ensure that histogram is working",
				},
				Buckets: nil,
				Labels:  nil,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			txn := NewTransaction(tt.opts)
			txn.Start(tt.values...)
			time.Sleep(10 * time.Millisecond)
			txn.End()

			labels := ""
			for i := 0; i < len(tt.opts.Labels); i++ {
				labels += fmt.Sprintf(`%s="%s"`, tt.opts.Labels[i], tt.values[i])
				if i < len(tt.opts.Labels)-1 {
					labels += ","
				}
			}

			fqName := buildFQName(tt.opts.MetaOpts)

			metrics := dump(prometheus.DefaultGatherer)
			assert.Contains(t, metrics, tt.opts.MetaOpts.Help)
			if len(tt.opts.Labels) > 0 {
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="0.005"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="0.01"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="0.025"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="0.05"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="0.1"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="0.25"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="0.5"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="1"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="2.5"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="5"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="10"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="+Inf"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_count{%s}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_sum{%s}`, fqName, labels))
			} else {
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="0.005"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="0.01"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="0.025"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="0.05"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="0.1"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="0.25"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="0.5"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="1"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="2.5"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="5"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="10"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="+Inf"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_count`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_sum`, fqName))
			}
		})
	}
}

func TestHistogramObserve(t *testing.T) {
	tests := []struct {
		name   string
		opts   *HistogramOpts
		values []string
	}{
		{
			name: "WithLabels",
			opts: &HistogramOpts{
				MetaOpts: MetaOpts{
					Namespace: "test",
					Subsystem: "unit",
					Name:      "histogram_with_labels",
					Help:      "ensure that histogram vec is working",
				},
				Buckets: DefBuckets,
				Labels:  []string{"t1"},
			},
			values: []string{"v1"},
		},
		{
			name: "WithoutLabels",
			opts: &HistogramOpts{
				MetaOpts: MetaOpts{
					Namespace: "test",
					Subsystem: "unit",
					Name:      "histogram_without_labels",
					Help:      "ensure that histogram is working",
				},
				Buckets: nil,
				Labels:  nil,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			txn := NewTransaction(tt.opts)
			txn.Observe(0.1, tt.values...)

			labels := ""
			for i := 0; i < len(tt.opts.Labels); i++ {
				labels += fmt.Sprintf(`%s="%s"`, tt.opts.Labels[i], tt.values[i])
				if i < len(tt.opts.Labels)-1 {
					labels += ","
				}
			}

			fqName := buildFQName(tt.opts.MetaOpts)

			metrics := dump(prometheus.DefaultGatherer)
			assert.Contains(t, metrics, tt.opts.MetaOpts.Help)
			if len(tt.opts.Labels) > 0 {
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="0.005"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="0.01"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="0.025"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="0.05"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="0.1"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="0.25"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="0.5"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="1"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="2.5"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="5"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="10"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{%s,le="+Inf"}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_count{%s}`, fqName, labels))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_sum{%s}`, fqName, labels))
			} else {
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="0.005"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="0.01"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="0.025"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="0.05"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="0.1"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="0.25"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="0.5"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="1"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="2.5"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="5"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="10"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_bucket{le="+Inf"}`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_count`, fqName))
				assert.Contains(t, metrics, fmt.Sprintf(`%s_sum`, fqName))
			}
		})
	}
}
