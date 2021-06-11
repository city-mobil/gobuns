package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/city-mobil/gobuns/zlog"
	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	Lag() int64
	ReadLag(context.Context) (int64, error)

	Offset() int64
	SetOffset(int64) error
	SetOffsetAt(context.Context, time.Time) error

	CommitMessages(context.Context, ...kafka.Message) error
	FetchMessage(context.Context) (kafka.Message, error)

	Ping() error
	Name() string
	ComponentType() string
	Close() error
}

type consumer struct {
	logger     zlog.Logger
	config     *ConsumerConfig
	reader     *kafka.Reader
	stats      *consumerStats
	onceCloser sync.Once
	stop       chan struct{}
}

func (c *consumer) Lag() int64 {
	return c.reader.Lag()
}

func (c *consumer) ReadLag(ctx context.Context) (int64, error) {
	return c.reader.ReadLag(ctx)
}

func (c *consumer) Offset() int64 {
	return c.reader.Offset()
}

func (c *consumer) SetOffset(offset int64) error {
	return c.reader.SetOffset(offset)
}

func (c *consumer) SetOffsetAt(ctx context.Context, t time.Time) error {
	return c.reader.SetOffsetAt(ctx, t)
}

func (c *consumer) CommitMessages(ctx context.Context, messages ...kafka.Message) error {
	return c.reader.CommitMessages(ctx, messages...)
}

func (c *consumer) FetchMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.FetchMessage(ctx)
}

func (c *consumer) Ping() error {
	return nil
}

func (c *consumer) Name() string {
	return "consumer"
}

func (c *consumer) ComponentType() string {
	return defaultType
}

func (c *consumer) Close() error {
	return c.reader.Close()
}

func NewConsumer(
	logger zlog.Logger,
	cfg *ConsumerConfig,
) Consumer {
	rdr := cfg.toKafkaReader(logger)

	c := &consumer{
		logger: logger,
		config: cfg,
		reader: rdr,
		stats:  newConsumerStats(),
		stop:   make(chan struct{}),
	}

	updater := &statsUpdater{
		refreshInterval: cfg.StatsConfig.RefreshInterval,
		stop:            c.stop,
		enabled:         cfg.StatsConfig.Enabled,
	}
	updater.run(func() {
		st := rdr.Stats()
		c.stats.updateFromStats(&st)
	})

	return c
}
