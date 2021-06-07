package rabbit

import (
	"time"

	"github.com/city-mobil/gobuns/rabbit/metrics"
	"github.com/city-mobil/gobuns/zlog"
	"github.com/streadway/amqp"
	"go.uber.org/atomic"
)

type Consumer interface {
	Consume()
	Close()
}

type consumer struct {
	logger          zlog.Logger
	conn            Connector
	cfg             *ConsumerCfg
	receiver        chan amqp.Delivery
	notifyClose     chan struct{}
	activeConsumers *atomic.Int64
}

func NewConsumer(logger zlog.Logger, cfg *ConsumerCfg, conn Connector, receiver chan amqp.Delivery) Consumer {
	cfg.withDefaults()
	return &consumer{
		logger: logger.With().
			Str("consumer", cfg.Consumer).
			Str("queue", cfg.Queue).
			Logger(),
		cfg:             cfg,
		conn:            conn,
		notifyClose:     make(chan struct{}),
		receiver:        receiver,
		activeConsumers: atomic.NewInt64(0),
	}
}

func (c *consumer) Consume() {
	go c.consume()
}

func (c *consumer) Close() {
	c.notifyClose <- struct{}{}
}

func (c *consumer) consume() {
	connect := func() {
		ch, err := c.conn.GetChannel()
		if err != nil {
			c.logger.Err(err).Msg("failed to get connection")
			return
		}

		delivery, err := ch.Consume(
			c.cfg.Queue,
			c.cfg.Consumer,
			c.cfg.AutoAck,
			c.cfg.Exclusive,
			false,
			c.cfg.NoWait,
			c.cfg.Args,
		)
		if err != nil {
			c.logger.Err(err).Msg("failed to consume queue")
			return
		}
		c.runConsumers(delivery)
	}
	connect()

	heartbeatTicker := time.NewTicker(c.cfg.Heartbeat)
	reloadConsumers := connect
	for {
		select {
		case <-heartbeatTicker.C:
			if !c.checkConsumers() {
				c.logger.Info().Msg("reloading consumers")
				reloadConsumers()
			}
		case <-c.notifyClose:
			c.logger.Info().Msg("consumer closed")
			return
		}
	}
}

func (c *consumer) runConsumers(delivery <-chan amqp.Delivery) {
	receiveMsg := func(msg amqp.Delivery) {
		c.receiver <- msg
		if !c.cfg.AutoAck {
			err := msg.Ack(false)
			if err != nil {
				c.logger.Err(err).Msg("failed to ack message receiving")
				return
			}
		}
		metrics.MessageConsumed(c.cfg.Queue)
	}

	consume := func() {
		for mqMessage := range delivery {
			receiveMsg(mqMessage)
		}
		c.logger.Info().Msg("consumer was closed")
		c.activeConsumers.Dec()
	}

	var i int64
	consumersCount := c.cfg.ConsumersCount - c.activeConsumers.Load()
	for i = 0; i < consumersCount; i++ {
		go consume()
		c.activeConsumers.Inc()
	}
}

func (c *consumer) checkConsumers() bool {
	return c.cfg.ConsumersCount == c.activeConsumers.Load()
}
