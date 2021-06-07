package rabbit

import (
	"fmt"
	"time"

	"github.com/city-mobil/gobuns/config"

	"github.com/streadway/amqp"
)

const (
	defURI = "default_unknown_uri"

	defReconnectDelay    = 3 * time.Second
	defConsumerHeartbeat = 3 * time.Second
	defConnectTimeout    = 500 * time.Second
	defConsumerCount     = 3
	defName              = "default_unknown_name"
	defHeartbeat         = 60 * time.Second
)

type Config struct {
	Addr           string
	Login          string
	Password       string
	ReconnectDelay time.Duration
	ConnectTimeout time.Duration
	Heartbeat      time.Duration
}

func (c *Config) GetURI() string {
	return fmt.Sprintf("amqp://%s:%s@%s/", c.Login, c.Password, c.Addr)
}

func NewConnectorConfig(prefix string) func() *Config {
	n := func(opt string) string {
		return fmt.Sprintf("%s.%s", prefix, opt)
	}
	var (
		addr           = config.String(n("addr"), defURI, "RabbitMQ address")
		login          = config.String(n("login"), defURI, "RabbitMQ login")
		password       = config.String(n("password"), defURI, "RabbitMQ password")
		reconnectDelay = config.Duration(n("reconnect_delay"), defReconnectDelay, "Sleep between reconnections")
		connectTimeout = config.Duration(n("connect_timeout"), defConnectTimeout, "Connection timeout")
		heartbeat      = config.Duration(n("heartbeat"), defHeartbeat, "Server heartbeat interval")
	)

	return func() *Config {
		return &Config{
			Addr:           *addr,
			Login:          *login,
			Password:       *password,
			ReconnectDelay: *reconnectDelay,
			ConnectTimeout: *connectTimeout,
			Heartbeat:      *heartbeat,
		}
	}
}

type ConsumerCfg struct {
	Queue          string
	ConsumersCount int64
	Consumer       string
	AutoAck        bool
	Exclusive      bool
	NoWait         bool
	Args           amqp.Table
	Heartbeat      time.Duration
}

func (c *ConsumerCfg) withDefaults() {
	if c.ConsumersCount == 0 {
		c.ConsumersCount = defConsumerCount
	}
	if c.Heartbeat == 0 {
		c.Heartbeat = defConsumerHeartbeat
	}
	if c.Consumer == "" {
		c.Consumer = defName
	}
}

func NewConsumerConfig(prefix string) func() *ConsumerCfg {
	n := func(opt string) string {
		return fmt.Sprintf("%s.%s", prefix, opt)
	}
	var (
		queue          = config.String(n("queue_name"), defName, "Consuming queue name")
		consumerName   = config.String(n("consumer"), defName, "Consumer name")
		consumersCount = config.Int64(n("consumers_count"), defConsumerCount, "Count of consumer workers")
		heartbeat      = config.Duration(n("heartbeat"), defConsumerHeartbeat, "Interval between checking count of consumers")
		autoAck        = config.Bool(n("auto_ack"), false, "Server will acknowledge deliveries to consumer prior to writing delivery to network")
		exclusive      = config.Bool(n("exclusive"), false, "Server will ensure this is the sole consumer from this queue")
		noWait         = config.Bool(n("no_wait"), false, "Don't wait for server to confirm request and immediately begin deliveries")
	)

	return func() *ConsumerCfg {
		return &ConsumerCfg{
			Queue:          *queue,
			ConsumersCount: *consumersCount,
			Consumer:       *consumerName,
			AutoAck:        *autoAck,
			Exclusive:      *exclusive,
			NoWait:         *noWait,
			Heartbeat:      *heartbeat,
		}
	}
}
