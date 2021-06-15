package kafka

import (
	"github.com/city-mobil/gobuns/promlib"
	"github.com/segmentio/kafka-go"
)

var (
	additionalLabels = []string{
		"topic",
	}
)

type consumerStats struct {
	lag           promlib.GaugeVec
	bytes         *promlib.Event
	errors        *promlib.Event
	offset        promlib.GaugeVec
	dialTimeAvg   promlib.GaugeVec
	readTimeAvg   promlib.GaugeVec
	fetchSize     *promlib.Event
	fetchBytesAvg promlib.GaugeVec
	messages      *promlib.Event
	fetches       *promlib.Event
	dials         *promlib.Event
}

func newConsumerStats() *consumerStats {
	prefix := "kafka_consumer_"

	return &consumerStats{
		lag: promlib.NewGaugeVec(promlib.GaugeOptions{
			MetaOpts: promlib.MetaOpts{
				Name:      prefix + "lag",
				Namespace: promlib.GetGlobalNamespace(),
				Help:      "Kafka consumer current lag",
			},
		}, additionalLabels),
		bytes: &promlib.Event{
			Name:      prefix + "bytes_count",
			Namespace: promlib.GetGlobalNamespace(),
			Help:      "Kafka consumer bytes sent count",
		},
		errors: &promlib.Event{
			Name:      prefix + "errors_count",
			Namespace: promlib.GetGlobalNamespace(),
			Help:      "Kafka consumer errors count",
		},
		offset: promlib.NewGaugeVec(promlib.GaugeOptions{
			MetaOpts: promlib.MetaOpts{
				Namespace: promlib.GetGlobalNamespace(),
				Name:      prefix + "offset",
				Help:      "Kafka consumer offset",
			},
		}, additionalLabels),
		dialTimeAvg: promlib.NewGaugeVec(promlib.GaugeOptions{
			MetaOpts: promlib.MetaOpts{
				Name:      prefix + "dial_time_avg",
				Namespace: promlib.GetGlobalNamespace(),
				Help:      "Kafka consumer dial time",
			},
		}, additionalLabels),
		readTimeAvg: promlib.NewGaugeVec(promlib.GaugeOptions{
			MetaOpts: promlib.MetaOpts{
				Name:      prefix + "read_time_avg",
				Namespace: promlib.GetGlobalNamespace(),
				Help:      "Kafka producer read time",
			},
		}, additionalLabels),
		fetchBytesAvg: promlib.NewGaugeVec(promlib.GaugeOptions{
			MetaOpts: promlib.MetaOpts{
				Namespace: promlib.GetGlobalNamespace(),
				Name:      prefix + "fetch_bytes_avg",
				Help:      "Kafka producer fetch bytes average",
			},
		}, additionalLabels),
		messages: &promlib.Event{
			Name:      prefix + "messages_count",
			Namespace: promlib.GetGlobalNamespace(),
			Help:      "Kafka consumer errors count",
		},
		fetches: &promlib.Event{
			Name:      prefix + "fetches_count",
			Namespace: promlib.GetGlobalNamespace(),
			Help:      "Kafka consumer errors count",
		},
		dials: &promlib.Event{
			Name:      prefix + "dials_count",
			Namespace: promlib.GetGlobalNamespace(),
			Help:      "Kafka consumer errors count",
		},
	}
}

func (c *consumerStats) updateFromStats(st *kafka.ReaderStats) {
	if st == nil {
		return
	}

	topicLabels := promlib.Labels{
		"topic": st.Topic,
	}

	c.lag.Set(float64(st.Lag), st.Topic)
	promlib.AddCntEventWithLabels(c.bytes, topicLabels, float64(st.Bytes))
	promlib.AddCntEventWithLabels(c.errors, topicLabels, float64(st.Errors))
	c.offset.Set(float64(st.Offset), st.Topic)
	c.dialTimeAvg.Set(float64(st.DialTime.Avg), st.Topic)
	c.readTimeAvg.Set(float64(st.DialTime.Avg), st.Topic)
	c.fetchBytesAvg.Set(float64(st.FetchBytes.Avg), st.Topic)
	promlib.AddCntEventWithLabels(c.messages, topicLabels, float64(st.Messages))
	promlib.AddCntEventWithLabels(c.fetches, topicLabels, float64(st.Fetches))
	promlib.AddCntEventWithLabels(c.dials, topicLabels, float64(st.Dials))
}
