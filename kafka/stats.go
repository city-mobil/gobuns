package kafka

import (
	"time"

	"github.com/city-mobil/gobuns/config"
	"github.com/city-mobil/gobuns/promlib"
	"github.com/segmentio/kafka-go"
)

const (
	defaultStatsCollectionEnabled = true
	defaultStatsPrefix            = "kafka"
	defaultStatsRefreshInterval   = time.Second
)

type StatsConfig struct {
	StatsPrefix     string
	RefreshInterval time.Duration
	Enabled         bool
}

func newStatsConfig(prefix string) func() *StatsConfig {
	var (
		statsPrefix     = config.String(prefix+"stats.prefix", defaultStatsPrefix, "Kafka producer stats prefix.")
		enabled         = config.Bool(prefix+"stats.enabled", defaultStatsCollectionEnabled, "Kafka producer stats enabled.")
		refreshInterval = config.Duration(prefix+"stats.refresh_interval", defaultStatsRefreshInterval, "Kafka producer stats refresh interval")
	)

	return func() *StatsConfig {
		return &StatsConfig{
			StatsPrefix:     *statsPrefix,
			RefreshInterval: *refreshInterval,
			Enabled:         *enabled,
		}
	}
}

type stats struct {
	writesCounter   *promlib.Event
	messagesCounter *promlib.Event
	bytesCounter    *promlib.Event
	errorsCounter   *promlib.Event
	batchTimeAvg    promlib.Gauge
	batchTimeMax    promlib.Gauge
}

func (s *stats) updateFromStats(st *kafka.WriterStats) {
	promlib.AddCntEvent(s.writesCounter, float64(st.Writes))
	promlib.AddCntEvent(s.messagesCounter, float64(st.Messages))
	promlib.AddCntEvent(s.bytesCounter, float64(st.Bytes))
	promlib.AddCntEvent(s.errorsCounter, float64(st.Errors))
	s.batchTimeAvg.Set(float64(st.BatchTime.Avg.Milliseconds()))
	s.batchTimeMax.Set(float64(st.BatchTime.Max.Milliseconds()))
}

func newProducerStats(prefix string) *stats {
	if prefix == "" {
		prefix = "kafka_producer_"
	} else {
		prefix += "_producer_"
	}

	return &stats{
		writesCounter: &promlib.Event{
			Name:      prefix + "count_writes",
			Namespace: promlib.GetGlobalNamespace(),
			Help:      "Kafka producer writes count",
		},
		messagesCounter: &promlib.Event{
			Name:      prefix + "count_messages",
			Namespace: promlib.GetGlobalNamespace(),
			Help:      "Kafka producer messages sent count",
		},
		bytesCounter: &promlib.Event{
			Name:      prefix + "count_bytes",
			Namespace: promlib.GetGlobalNamespace(),
			Help:      "Kafka producer bytes sent count",
		},
		errorsCounter: &promlib.Event{
			Name:      prefix + "count_errors",
			Namespace: promlib.GetGlobalNamespace(),
			Help:      "Kafka producer errors count",
		},
		batchTimeAvg: promlib.NewGauge(
			promlib.GaugeOptions{
				MetaOpts: promlib.MetaOpts{
					Namespace: promlib.GetGlobalNamespace(),
					Name:      prefix + "batch_time_avg",
					Help:      "Kafka producer maximum batch time",
				},
			}),
		batchTimeMax: promlib.NewGauge(
			promlib.GaugeOptions{
				MetaOpts: promlib.MetaOpts{
					Namespace: promlib.GetGlobalNamespace(),
					Name:      prefix + "batch_time_max",
					Help:      "Kafka producer maximum batch time",
				},
			},
		),
	}
}
