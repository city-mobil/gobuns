package promlib

import "github.com/prometheus/client_golang/prometheus"

type GaugeVec struct {
	gaugeMetric *prometheus.GaugeVec
}

func NewGaugeVec(opts GaugeOptions, labels []string) GaugeVec {
	gaugeMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: opts.Namespace,
		Subsystem: opts.Subsystem,
		Name:      opts.Name,
		Help:      opts.Help,
	}, labels)
	prometheus.MustRegister(gaugeMetric)

	return GaugeVec{
		gaugeMetric: gaugeMetric,
	}
}

func (g GaugeVec) Set(val float64, labelValues ...string) {
	g.gaugeMetric.WithLabelValues(labelValues...).Set(val)
}

func (g GaugeVec) Inc(labelValues ...string) {
	g.gaugeMetric.WithLabelValues(labelValues...).Inc()
}

func (g GaugeVec) Dec(labelValues ...string) {
	g.gaugeMetric.WithLabelValues(labelValues...).Dec()
}

func (g GaugeVec) Add(val float64, labelValues ...string) {
	g.gaugeMetric.WithLabelValues(labelValues...).Add(val)
}

func (g GaugeVec) Sub(val float64, labelValues ...string) {
	g.gaugeMetric.WithLabelValues(labelValues...).Sub(val)
}
