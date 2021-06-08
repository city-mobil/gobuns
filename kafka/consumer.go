package kafka

import "github.com/city-mobil/gobuns/zlog"

type Consumer interface {
}

type consumer struct {
}

func NewConsumer(
	logger zlog.Logger,
	cfg *ConsumerConfig,
) Consumer {
	rdr := cfg.toKafkaReader()
	rdr.Stats()
	return nil
}
