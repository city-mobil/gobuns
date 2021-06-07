package kafka

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/city-mobil/gobuns/barber"
	"github.com/city-mobil/gobuns/zlog"
	"github.com/segmentio/kafka-go"
)

const (
	breakerServerID = 1
)

const (
	defaultName = "producer"
	defaultType = "kafka"
)

var (
	ErrBrokerUnavailable = errors.New("kafka broker is unavailable")
)

// Producer is a kafka producer interface.
type Producer interface {
	// Produce produces given messages to kafka.
	Produce(context.Context, ...kafka.Message) error

	// SetCompletionCallback sets completion callback for asynchronous producer.
	//
	// For synchronous producer this callback is never called.
	SetCompletionCallback(CompletionCallback)

	Ping() error
	Name() string
	ComponentType() string
	ComponentID() string

	// Close closes underlying producer.
	// Close flushes pending writes, and waits for all writes to complete before
	// returning. Calling Close also prevents new writes from being submitted to
	// the writer, further calls to WriteMessages and the like will fail with
	// io.ErrClosedPipe.
	Close() error
}

type producer struct {
	logger                zlog.Logger
	config                *ProducerConfig
	writer                *kafka.Writer
	breaker               barber.Barber
	stats                 *stats
	completionCallbackSet uint32

	onceCloser sync.Once
	stop       chan struct{}
}

// Produce produces given messages to kafka.
func (p *producer) Produce(ctx context.Context, messages ...kafka.Message) error {
	if p.writer.Async && atomic.LoadUint32(&p.completionCallbackSet) == 0 {
		p.logger.Warn().Msgf("CompletionCallback is not set. We highly recommend to do this. See README")
	}

	if p.config.CircuitBreakerEnabled {
		alive := p.breaker.IsAvailable(breakerServerID, time.Now())
		if !alive {
			return ErrBrokerUnavailable
		}
	}

	err := p.writer.WriteMessages(ctx, messages...)
	if err != nil {
		p.breaker.AddError(breakerServerID, time.Now())
	}

	return err
}

// SetCompletionCallback sets completion callback for asynchronous producer.
//
// For synchronous producer this callback is never called.
func (p *producer) SetCompletionCallback(cb func(messages []kafka.Message, err error)) {
	atomic.StoreUint32(&p.completionCallbackSet, 1)

	if p.writer.Async {
		// For async producer it is important to set our error wrapper
		// otherwise we would not update internal circuit breaker state.
		p.writer.Completion = newOnErrorCallbackWrapper(p.breaker, cb)
	} else {
		p.writer.Completion = cb
	}
}

func (p *producer) Ping() error {
	alive := p.breaker.IsAvailable(breakerServerID, time.Now())
	if !alive {
		return ErrBrokerUnavailable
	}

	return nil
}

func (p *producer) Name() string {
	return defaultName
}

func (p *producer) ComponentType() string {
	return defaultType
}

func (p *producer) ComponentID() string {
	return p.writer.Addr.String()
}

// Close closes underlying producer.
// Close flushes pending writes, and waits for all writes to complete before
// returning. Calling Close also prevents new writes from being submitted to
// the writer, further calls to WriteMessages and the like will fail with
// io.ErrClosedPipe.
func (p *producer) Close() error {
	p.onceCloser.Do(func() {
		close(p.stop)
	})
	return p.writer.Close()
}

func (p *producer) runStatsUpdater() {
	if !p.config.StatsConfig.Enabled {
		return
	}

	go func() {
		t := time.NewTicker(p.config.StatsConfig.RefreshInterval)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				st := p.writer.Stats()
				p.stats.updateFromStats(&st)

			case <-p.stop:
				return
			}
		}
	}()
}

func newWriter(logger zlog.Logger, cfg *ProducerConfig) *kafka.Writer {
	writer := cfg.toKafkaWriter()

	// NOTE(a.petrukhin): minor optimization. Logger is not initialized if the
	// level is zlog.Disabled
	if cfg.LogLevel != zlog.Disabled {
		writer.Logger = kafka.LoggerFunc(func(fmt string, args ...interface{}) {
			logger.Level(cfg.LogLevel).Log().Msgf(fmt, args...)
		})
	}
	if cfg.ErrorLogLevel != zlog.Disabled {
		writer.ErrorLogger = kafka.LoggerFunc(func(fmt string, args ...interface{}) {
			logger.Error().Msgf(fmt, args...)
		})
	}

	return writer
}

func newOnErrorCallbackWrapper(breaker barber.Barber, next CompletionCallback) CompletionCallback {
	cb := &completionCallbackOnError{
		onErr: func(_ error) {
			breaker.AddError(breakerServerID, time.Now())
		},
		next: next,
	}

	return cb.exec
}

// NewAsyncProducer creates new asynchronous producer with no completion callback.
//
// In asynchronous producer user does not wait for the acknowledge from the kafka broker.
// If any error occurs, the error can be found in CompletionCallback.
// For asynchronous producer CompletionCallback MUST BE SET in order to have possibility to retry messages.
func NewAsyncProducer(
	logger zlog.Logger,
	cfg *ProducerConfig,
) Producer {
	breaker := barber.NewBarber([]int{breakerServerID}, cfg.CircuitBreakerConfig)

	writer := newWriter(logger, cfg)
	writer.Async = true
	writer.Completion = newOnErrorCallbackWrapper(breaker, CompletionCallbackDiscard)

	p := &producer{
		logger:  logger,
		config:  cfg,
		writer:  writer,
		stats:   newProducerStats(cfg.StatsConfig.StatsPrefix),
		stop:    make(chan struct{}),
		breaker: breaker,
	}
	p.runStatsUpdater()

	return p
}

// NewSyncProducer creates new synchronous producer.
//
// In synchronous mode user waits for acknowledge of each produced message (or a group of messages if
// multiple messages are produced at the same time.)
func NewSyncProducer(logger zlog.Logger, cfg *ProducerConfig) Producer {
	writer := newWriter(logger, cfg)
	writer.Async = false

	p := &producer{
		logger:  logger,
		config:  cfg,
		writer:  writer,
		stats:   newProducerStats(cfg.StatsConfig.StatsPrefix),
		stop:    make(chan struct{}),
		breaker: barber.NewBarber([]int{breakerServerID}, cfg.CircuitBreakerConfig),
	}
	p.runStatsUpdater()

	return p
}
