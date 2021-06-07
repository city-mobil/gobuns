package kafka

import (
	"net"
	"time"

	"github.com/city-mobil/gobuns/barber"
	"github.com/city-mobil/gobuns/config"
	"github.com/city-mobil/gobuns/zlog"

	"github.com/segmentio/kafka-go"
)

const (
	defaultMaxRetries            = 3
	defaultQueueMaxMessages      = 10000
	defaultQueueBufferMaxSize    = 1048576 // 1 MiB
	defaultQueueBufferingTimeout = 20 * time.Millisecond
	defaultReadTimeout           = 3 * time.Second
	defaultWriteTimeout          = 3 * time.Second
	defaultDialTimeout           = 3 * time.Second
	defaultCompression           = 0
	defaultLogLevel              = zlog.Disabled
	defaultCircuitBreakerEnabled = false
)

const (
	BalancerRoundRobin = "roundrobin"
	BalancerCRC32      = "crc32"
	BalancerMurMur2    = "murmur2"
)

type ProducerConfig struct { //nolint:maligned
	Brokers               []string
	QueueMaxBytesSize     int64
	QueueBufferingTimeout time.Duration
	QueueMaxMessages      int
	MaxRetries            int
	ReadTimeout           time.Duration
	WriteTimeout          time.Duration
	DialTimeout           time.Duration
	RequiredAcks          kafka.RequiredAcks
	Balancer              kafka.Balancer
	LogLevel              zlog.Level
	ErrorLogLevel         zlog.Level
	Compression           kafka.Compression
	StatsConfig           *StatsConfig
	CircuitBreakerConfig  *barber.Config
	CircuitBreakerEnabled bool
}

func NewProducerConfig(prefix string) func() *ProducerConfig {
	if prefix != "" {
		prefix += ".kafka."
	} else {
		prefix = "kafka."
	}

	o := func(opt string) string {
		return prefix + opt
	}

	var (
		brokers = config.StringSlice(
			o("producer.brokers"),
			[]string{},
			"Kafka broker addresses.",
		)
		balancer = config.String(
			o("producer.balancer"),
			BalancerRoundRobin,
			"Kafka producer balancer",
		)
		maxRetries = config.Int(
			o("producer.max_retries"),
			defaultMaxRetries,
			"Kafka message max retries. Analog for "+
				"messages.send.max.retries rdkafka option.",
		)
		maxMessages = config.Int(
			o("producer.queue.max_messages"),
			defaultQueueMaxMessages,
			"Kafka producer internal "+
				"queue max messages. Analog for queue.buffering.max.messages rdkafka option.",
		)
		maxQueueSize = config.Int64(
			o("producer.queue.max_bytes"),
			defaultQueueBufferMaxSize,
			"Kafka producer "+
				"internal queue max size in bytes. Analog for queue.buffering.max.kbytes rdkafka option.",
		)
		queueBufferingTimeout = config.Duration(
			o("producer.queue.buffering_timeout"),
			defaultQueueBufferingTimeout,
			"Kafka producer internal queue max buffering timeout. Analog for queue.buffering.max.ms rdkafka"+
				"option.",
		)
		readTimeout = config.Duration(
			o("producer.net.read_timeout"),
			defaultReadTimeout,
			"Kafka producer network read timeout.")
		writeTimeout = config.Duration(
			o("producer.net.write_timeout"),
			defaultWriteTimeout,
			"Kafka producer network write timeout.",
		)
		dialTimeout = config.Duration(
			o("producer.net.dial_timeout"),
			defaultDialTimeout,
			"Kafka producer network dial timeout.",
		)
		requiredAcks = config.Int(
			o("producer.required_acks"),
			0,
			"Kafka producer required acks. Analog for request.required.acks rdkafka option.",
		)
		compression = config.Int8(
			o("producer.compression"),
			defaultCompression,
			"Kafka producer compression. Analog for compression.codec rdkafka option.",
		)
		logLevel = config.Int8(
			o("log.level"),
			int8(defaultLogLevel),
			"Kafka producer logger log level.",
		)
		errorLogLevel = config.Int8(
			o("log.errors_level"),
			int8(defaultLogLevel),
			"Kafka producer error-logger log level",
		)
		statsConfig    = newStatsConfig(prefix)
		breakerEnabled = config.Bool(o("breaker.enabled"), defaultCircuitBreakerEnabled, "Kafka circuit breaker mode.")
		breakerConfig  = barber.NewConfig(o("breaker"))
	)

	return func() *ProducerConfig {
		return &ProducerConfig{
			Brokers:               *brokers,
			Balancer:              balancerToKafkaBalancer(*balancer),
			MaxRetries:            *maxRetries,
			QueueMaxMessages:      *maxMessages,
			QueueMaxBytesSize:     *maxQueueSize,
			QueueBufferingTimeout: *queueBufferingTimeout,
			ReadTimeout:           *readTimeout,
			WriteTimeout:          *writeTimeout,
			DialTimeout:           *dialTimeout,
			RequiredAcks:          kafka.RequiredAcks(*requiredAcks),
			Compression:           kafka.Compression(*compression),
			LogLevel:              zlog.Level(*logLevel),
			ErrorLogLevel:         zlog.Level(*errorLogLevel),
			StatsConfig:           statsConfig(),
			CircuitBreakerEnabled: *breakerEnabled,
			CircuitBreakerConfig:  breakerConfig(),
		}
	}
}

func (c *ProducerConfig) toKafkaWriter() *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(c.Brokers...),
		Balancer:     c.Balancer,
		MaxAttempts:  c.MaxRetries,
		BatchSize:    c.QueueMaxMessages,
		BatchBytes:   c.QueueMaxBytesSize,
		BatchTimeout: c.QueueBufferingTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		RequiredAcks: c.RequiredAcks,
		Async:        false,
		Completion:   nil,
		Compression:  c.Compression,
		Transport: &kafka.Transport{
			Dial: (&net.Dialer{
				// NOTE(a.petrukhin): DualStack is enabled by default.
				// See RFC 6555 for further information
				Timeout: c.DialTimeout,
			}).DialContext,
		},
	}
}
