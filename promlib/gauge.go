package promlib

import "github.com/prometheus/client_golang/prometheus"

type GaugeOptions struct {
	MetaOpts
}

type Gauge struct {
	gaugeMetric prometheus.Gauge
}

func NewGauge(opts GaugeOptions) Gauge {
	gaugeMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   opts.Namespace,
		Subsystem:   opts.Subsystem,
		Name:        opts.Name,
		Help:        opts.Help,
		ConstLabels: opts.ConstLabels,
	})
	prometheus.MustRegister(gaugeMetric)

	return Gauge{
		gaugeMetric: gaugeMetric,
	}
}

func (g Gauge) Set(val float64) {
	g.gaugeMetric.Set(val)
}

func (g Gauge) Inc() {
	g.gaugeMetric.Inc()
}

func (g Gauge) Dec() {
	g.gaugeMetric.Dec()
}

func (g Gauge) Add(val float64) {
	g.gaugeMetric.Add(val)
}

func (g Gauge) Sub(val float64) {
	g.gaugeMetric.Sub(val)
}
