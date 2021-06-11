package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/city-mobil/gobuns/zlog"
	"github.com/segmentio/kafka-go"
)

// Consumer is a kafka consumer interface.
type Consumer interface {
	// Lag returns the lag of the last message returned by ReadMessage, or -1
	// if r is backed by a consumer group.
	Lag() int64

	// ReadLag returns the current lag of the reader by fetching the last offset of
	// the topic and partition and computing the difference between that value and
	// the offset of the last message returned by ReadMessage.
	//
	// This method is intended to be used in cases where a program may be unable to
	// call ReadMessage to update the value returned by Lag, but still needs to get
	// an up to date estimation of how far behind the reader is. For example when
	// the consumer is not ready to process the next message.
	//
	// The function returns a lag of zero when the reader's current offset is
	// negative.
	ReadLag(context.Context) (int64, error)

	// Offset returns the current absolute offset of the reader, or -1
	// if r is backed by a consumer group.
	Offset() int64

	// SetOffset changes the offset from which the next batch of messages will be
	// read. The method fails with io.ErrClosedPipe if the reader has already been closed.
	//
	// From version 0.2.0, FirstOffset and LastOffset can be used to indicate the first
	// or last available offset in the partition. Please note while -1 and -2 were accepted
	// to indicate the first or last offset in previous versions, the meanings of the numbers
	// were swapped in 0.2.0 to match the meanings in other libraries and the Kafka protocol
	// specification.
	SetOffset(int64) error

	// SetOffsetAt changes the offset from which the next batch of messages will be
	// read given the timestamp t.
	//
	// The method fails if the unable to connect partition leader, or unable to read the offset
	// given the ts, or if the reader has been closed.
	SetOffsetAt(context.Context, time.Time) error

	// CommitMessages commits the list of messages passed as argument. The program
	// may pass a context to asynchronously cancel the commit operation when it was
	// configured to be blocking.
	CommitMessages(context.Context, ...kafka.Message) error

	// ReadMessage reads and return the next message from the r. The method call
	// blocks until a message becomes available, or an error occurs. The program
	// may also specify a context to asynchronously cancel the blocking operation.
	//
	// The method returns io.EOF to indicate that the reader has been closed.
	//
	// If consumer groups are used, ReadMessage will automatically commit the
	// offset when called. Note that this could result in an offset being committed
	// before the message is fully processed.
	//
	// If more fine grained control of when offsets are  committed is required, it
	// is recommended to use FetchMessage with CommitMessages instead.
	ReadMessage(context.Context) (kafka.Message, error)

	// FetchMessage reads and return the next message from the r. The method call
	// blocks until a message becomes available, or an error occurs. The program
	// may also specify a context to asynchronously cancel the blocking operation.
	//
	// The method returns io.EOF to indicate that the reader has been closed.
	//
	// FetchMessage does not commit offsets automatically when using consumer groups.
	// Use CommitMessages to commit the offset.
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

func (c *consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
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
