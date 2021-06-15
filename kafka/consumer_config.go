package kafka

import (
	"context"
	"net"
	"time"

	"github.com/city-mobil/gobuns/zlog"

	"github.com/city-mobil/gobuns/config"
	"github.com/segmentio/kafka-go"
)

const (
	GroupBalancerRange uint8 = iota
	GroupBalancerRoundRobin
)

type ConsumerDialerConfig struct {
	ClientID      string
	DialFunc      func(ctx context.Context, network string, addr string) (net.Conn, error)
	Timeout       time.Duration
	LocalAddr     string
	FallbackDelay time.Duration
	KeepAlive     time.Duration
}

type ConsumerConfig struct {
	GroupID                string
	Brokers                []string
	Topic                  string
	Partition              int
	DialerConfig           *ConsumerDialerConfig
	StatsConfig            *StatsConfig
	DialTimeout            time.Duration
	QueueCapacity          int
	MinBytes               int
	MaxBytes               int
	MaxWait                time.Duration
	ReadLagInterval        time.Duration
	GroupBalancers         []uint
	HeartBeatInterval      time.Duration
	CommitInterval         time.Duration
	PartitionWatchInterval time.Duration
	SessionTimeout         time.Duration
	RebalanceTimeout       time.Duration
	JoinGroupBackoff       time.Duration
	RetentionTime          time.Duration
	StartOffset            int64
	ReadBackoffMin         time.Duration
	ReadBackoffMax         time.Duration
	MaxAttempts            int
	IsolationLevel         int8
	WatchPartitionChanges  bool
	LogLevel               zlog.Level
	ErrorLogLevel          zlog.Level
}

func newDialerConfig(prefix string) func() *ConsumerDialerConfig {
	o := func(opt string) string {
		return prefix + "dialer." + opt
	}

	var (
		clientID      = config.String(o("client_id"), "", "Kafka Consumer Dialer ClientID.")
		timeout       = config.Duration(o("timeout"), 0, "Kafka Consumer Dialer Timeout.")
		localAddr     = config.String(o("local_addr"), "", "Kafka Consumer Dialer LocalAddr.")
		fallbackDelay = config.Duration(o("fallback_delay"), 0, "Kafka Consumer FallbackDelay.")
		keepAlive     = config.Duration("keep_alive", 0, "Kafka Consumer KeepAlive.")
	)

	return func() *ConsumerDialerConfig {
		return &ConsumerDialerConfig{
			ClientID:      *clientID,
			Timeout:       *timeout,
			LocalAddr:     *localAddr,
			FallbackDelay: *fallbackDelay,
			KeepAlive:     *keepAlive,
		}
	}
}

func NewConsumerConfig(prefix string) func() *ConsumerConfig {
	if prefix != "" {
		prefix += ".kafka.consumer."
	} else {
		prefix = "kafka.consumer."
	}

	const (
		consumerDefaultQueueCapacity   = 100
		consumerDefaultMaxWait         = 5 * time.Second
		consumerDefaultMinBytes        = 1
		consumerDefaultMaxBytes        = 0
		consumerDefaultReadLagInterval = time.Minute
	)

	o := func(opt string) string {
		return prefix + opt
	}

	var (
		groupID = config.String(
			o("consumer.group_id"),
			"",
			"Kafka consumer group_id.",
		)
		brokers = config.StringSlice(
			o("consumer.brokers"),
			[]string{},
			"Kafka consumer brokers.",
		)
		topic = config.String(
			o("consumer.topic"),
			"",
			"Kafka consumer topic.",
		)
		partition = config.Int(
			o("consumer.partition"),
			0,
			"Kafka Consumer Partition.",
		)
		dialerConfig = newDialerConfig(o("consumer.dialer"))
		dialTimeout  = config.Duration(
			o("consumer.net.dial_timeout"),
			defaultDialTimeout,
			"Kafka consumer dial timeout.",
		)
		queueCapacity = config.Int(
			o("consumer.queue_capacity"),
			consumerDefaultQueueCapacity,
			"Kafka consumer internal queue capacity",
		)
		minBytes = config.Int(
			o("consumer.fetch.min_bytes"),
			consumerDefaultMinBytes,
			"Kafka consumer min bytes to fetch on each request",
		)
		maxBytes = config.Int(
			o("consumer.fetch.max_bytes"),
			consumerDefaultMaxBytes,
			"Kafka consumer max bytes to fetch on each request.",
		)
		maxWait = config.Duration(
			o("consumer.max_wait"),
			consumerDefaultMaxWait,
			"Consumer Max Wait.",
		)
		readLagInterval = config.Duration(
			o("consumer.read_lag_interval"),
			consumerDefaultReadLagInterval,
			"Consumer frequency at which the reader lag is updated.",
		)
		groupBalancers = config.UintSlice(
			o("consumer.group_balancers"),
			[]uint{uint(GroupBalancerRange), uint(GroupBalancerRoundRobin)},
			"Consumer Group Balancers.",
		)
		heartBeatInterval = config.Duration(
			o("consumer.heart_beat_interval"),
			0,
			"Consumer Heartbeat Interval.",
		)
		commitInterval = config.Duration(
			o("consumer.commit_interval"),
			0,
			"Consumer Commit Interval.",
		)
		partitionWatchInterval = config.Duration(
			o("consumer.partition_watch_interval"),
			0,
			"Consumer Partition Watch Interval.",
		)
		watchPartitionChanges = config.Bool(
			o("watch_partition_changes"),
			false,
			"Kafka Consumer WatchPartitionChanges",
		)
		sessionTimeout = config.Duration(
			o("consumer.session_timeout"),
			0,
			"Consumer Session Timeout.",
		)
		rebalanceTimeout = config.Duration(
			o("consumer.rebalance_timeout"),
			0,
			"Consumer Rebalance Timeout.",
		)
		joinGroupBackoff = config.Duration(
			o("consumer.join_group_backoff"),
			0,
			"Consumer Join Group Backoff.",
		)
		retentionTime = config.Duration(
			o("consumer.retention_time"),
			0,
			"Consumer Retention Time.",
		)
		startOffset = config.Int64(
			o("consumer.start_offset"),
			0,
			"Consumer Start Offset.",
		)
		readBackoffMin = config.Duration(
			o("consumer.read_backoff_min"),
			0,
			"Consumer Read Backoff Min.",
		)
		readBackoffMax = config.Duration(
			o("consumer.read_backoff_max"),
			0,
			"Consumer Read Backoff Max.",
		)
		isolationLevel = config.Int8(
			o("consumer.isolation_level"),
			0,
			"Consumer Isolation Level.",
		)
		maxAttempts = config.Int(
			o("consumer.max_attempts"),
			0,
			"Consumer Max Attempts.",
		)
		logLevel = config.Int8(
			o("log.level"),
			int8(defaultLogLevel),
			"Kafka consumer logger log level.",
		)
		errorLogLevel = config.Int8(
			o("log.errors_level"),
			int8(defaultLogLevel),
			"Kafka consumer error-logger log level.",
		)
		statsConfig = newStatsConfig(o("consumer."))
	)

	return func() *ConsumerConfig {
		return &ConsumerConfig{
			GroupID:                *groupID,
			Brokers:                *brokers,
			Topic:                  *topic,
			Partition:              *partition,
			DialerConfig:           dialerConfig(),
			DialTimeout:            *dialTimeout,
			QueueCapacity:          *queueCapacity,
			MinBytes:               *minBytes,
			MaxBytes:               *maxBytes,
			MaxWait:                *maxWait,
			ReadLagInterval:        *readLagInterval,
			GroupBalancers:         *groupBalancers,
			HeartBeatInterval:      *heartBeatInterval,
			CommitInterval:         *commitInterval,
			PartitionWatchInterval: *partitionWatchInterval,
			WatchPartitionChanges:  *watchPartitionChanges,
			SessionTimeout:         *sessionTimeout,
			RebalanceTimeout:       *rebalanceTimeout,
			JoinGroupBackoff:       *joinGroupBackoff,
			RetentionTime:          *retentionTime,
			StartOffset:            *startOffset,
			ReadBackoffMin:         *readBackoffMin,
			ReadBackoffMax:         *readBackoffMax,
			IsolationLevel:         *isolationLevel,
			MaxAttempts:            *maxAttempts,
			LogLevel:               zlog.Level(*logLevel),
			ErrorLogLevel:          zlog.Level(*errorLogLevel),
			StatsConfig:            statsConfig(),
		}
	}
}

func (c *ConsumerConfig) toKafkaReader(logger zlog.Logger) *kafka.Reader {
	var groupBalancers []kafka.GroupBalancer
	if c.GroupBalancers != nil {
		for _, v := range c.GroupBalancers {
			switch uint8(v) {
			case GroupBalancerRoundRobin:
				groupBalancers = append(groupBalancers, kafka.RoundRobinGroupBalancer{})
			case GroupBalancerRange:
				groupBalancers = append(groupBalancers, kafka.RangeGroupBalancer{})
			}
		}
	}

	readerConfig := kafka.ReaderConfig{
		Brokers:   c.Brokers,
		GroupID:   c.GroupID,
		Topic:     c.Topic,
		Partition: c.Partition,
		Dialer: &kafka.Dialer{
			ClientID:        c.DialerConfig.ClientID, // TODO(a.petrukhin): think about it.
			DialFunc:        c.DialerConfig.DialFunc,
			Timeout:         c.DialerConfig.Timeout,
			LocalAddr:       kafka.TCP(c.DialerConfig.LocalAddr),
			DualStack:       true,
			FallbackDelay:   c.DialerConfig.FallbackDelay,
			KeepAlive:       c.DialerConfig.KeepAlive,
			Resolver:        nil,
			TLS:             nil,
			SASLMechanism:   nil,
			TransactionalID: "",
		},
		QueueCapacity:          c.QueueCapacity,
		MinBytes:               c.MinBytes,
		MaxBytes:               c.MaxBytes,
		MaxWait:                c.MaxWait,
		ReadLagInterval:        c.ReadLagInterval,
		GroupBalancers:         groupBalancers,
		HeartbeatInterval:      c.HeartBeatInterval,
		CommitInterval:         c.CommitInterval,
		PartitionWatchInterval: c.PartitionWatchInterval,
		WatchPartitionChanges:  c.WatchPartitionChanges,
		SessionTimeout:         c.SessionTimeout,
		RebalanceTimeout:       c.RebalanceTimeout,
		JoinGroupBackoff:       c.JoinGroupBackoff,
		RetentionTime:          c.RetentionTime,
		StartOffset:            c.StartOffset,
		ReadBackoffMin:         c.ReadBackoffMin,
		ReadBackoffMax:         c.ReadBackoffMax,
		Logger:                 nil,
		ErrorLogger:            nil,
		IsolationLevel:         kafka.IsolationLevel(c.IsolationLevel),
		MaxAttempts:            c.MaxAttempts,
	}

	if c.LogLevel != zlog.Disabled {
		readerConfig.Logger = kafka.LoggerFunc(func(fmt string, args ...interface{}) {
			logger.Level(c.LogLevel).Log().Msgf(fmt, args...)
		})
	}
	if c.ErrorLogLevel != zlog.Disabled {
		readerConfig.ErrorLogger = kafka.LoggerFunc(func(fmt string, args ...interface{}) {
			logger.Error().Msgf(fmt, args...)
		})
	}

	return kafka.NewReader(readerConfig)
}
