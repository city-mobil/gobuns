package kafka

import (
	"time"

	"github.com/city-mobil/gobuns/config"
	"github.com/city-mobil/gobuns/promlib"
	"github.com/segmentio/kafka-go"
)

const (
	defaultStatsCollectionEnabled = true
	defaultStatsRefreshInterval   = time.Second
)

type StatsConfig struct {
	RefreshInterval time.Duration
	Enabled         bool
}

func newStatsConfig(prefix string) func() *StatsConfig {
	var (
		enabled         = config.Bool(prefix+"stats.enabled", defaultStatsCollectionEnabled, "Kafka producer stats enabled.")
		refreshInterval = config.Duration(prefix+"stats.refresh_interval", defaultStatsRefreshInterval, "Kafka producer stats refresh interval")
	)

	return func() *StatsConfig {
		return &StatsConfig{
			RefreshInterval: *refreshInterval,
			Enabled:         *enabled,
		}
	}
}

type producerStats struct {
	writesCounter   *promlib.Event
	messagesCounter *promlib.Event
	bytesCounter    *promlib.Event
	errorsCounter   *promlib.Event
	batchTimeAvg    promlib.GaugeVec
	batchTimeMax    promlib.GaugeVec
}

func (s *producerStats) updateFromStats(st *kafka.WriterStats) {
	labels := promlib.Labels{
		"topic": st.Topic,
	}

	promlib.AddCntEventWithLabels(s.writesCounter, labels, float64(st.Writes))
	promlib.AddCntEventWithLabels(s.messagesCounter, labels, float64(st.Messages))
	promlib.AddCntEventWithLabels(s.bytesCounter, labels, float64(st.Bytes))
	promlib.AddCntEventWithLabels(s.errorsCounter, labels, float64(st.Errors))
	s.batchTimeAvg.Set(float64(st.BatchTime.Avg.Milliseconds()), st.Topic)
	s.batchTimeMax.Set(float64(st.BatchTime.Max.Milliseconds()), st.Topic)
}

func newProducerStats() *producerStats {
	prefix := "kafka_producer_"

	return &producerStats{
		writesCounter: &promlib.Event{
			Name:      prefix + "writes_count",
			Namespace: promlib.GetGlobalNamespace(),
			Help:      "Kafka producer writes count",
		},
		messagesCounter: &promlib.Event{
			Name:      prefix + "messages_count",
			Namespace: promlib.GetGlobalNamespace(),
			Help:      "Kafka producer messages sent count",
		},
		bytesCounter: &promlib.Event{
			Name:      prefix + "_count_bytes",
			Namespace: promlib.GetGlobalNamespace(),
			Help:      "Kafka producer bytes sent count",
		},
		errorsCounter: &promlib.Event{
			Name:      prefix + "_count_errors",
			Namespace: promlib.GetGlobalNamespace(),
			Help:      "Kafka producer errors count",
		},
		batchTimeAvg: promlib.NewGaugeVec(
			promlib.GaugeOptions{
				MetaOpts: promlib.MetaOpts{
					Namespace: promlib.GetGlobalNamespace(),
					Name:      prefix + "_batch_time_avg",
					Help:      "Kafka producer maximum batch time",
				},
			}, additionalLabels),
		batchTimeMax: promlib.NewGaugeVec(
			promlib.GaugeOptions{
				MetaOpts: promlib.MetaOpts{
					Namespace: promlib.GetGlobalNamespace(),
					Name:      prefix + "_batch_time_max",
					Help:      "Kafka producer maximum batch time",
				},
			}, additionalLabels,
		),
	}
}
