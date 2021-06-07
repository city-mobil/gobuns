package rabbit

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/city-mobil/gobuns/zlog"
	"github.com/streadway/amqp"
)

const (
	traceComponentName = "go-buns/rabbitMQ"
)

var (
	ErrChannelNotInitialized    = errors.New("channel is not initialized")
	ErrConnectionNotInitialized = errors.New("connection is not initialized")
)

const (
	defLocale = "en_US"
)

type Connector interface {
	GetChannel() (*amqp.Channel, error)
	GetConnection() (*amqp.Connection, error)
	GetHost() string
	Close() error
}

type connector struct {
	logger     zlog.Logger
	cfg        Config
	channel    *amqp.Channel
	connection *amqp.Connection
	mutex      *sync.RWMutex
}

func NewConnector(logger zlog.Logger, cfg Config) (Connector, error) {
	conn := &connector{
		logger: logger.With().
			Str("connector", "RabbitMQ").
			Str("addr", cfg.Addr).
			Logger(),
		cfg:   cfg,
		mutex: &sync.RWMutex{},
	}

	err := conn.connect()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *connector) GetHost() string {
	return c.cfg.Addr
}

func (c *connector) GetChannel() (*amqp.Channel, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.channel == nil {
		return nil, ErrChannelNotInitialized
	}

	return c.channel, nil
}

func (c *connector) GetConnection() (*amqp.Connection, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.connection == nil {
		return nil, ErrConnectionNotInitialized
	}

	return c.connection, nil
}

func (c *connector) Close() (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	defer func() {
		c.connection = nil
		c.channel = nil
	}()

	err = c.closeChannel()
	if err != nil {
		e := c.closeConnection()
		if e != nil {
			return fmt.Errorf("failed to close channel and connection %s %s", err, e)
		}
		return err
	}

	return c.closeConnection()
}

func (c *connector) closeConnection() error {
	if c.connection == nil {
		return nil
	}

	return c.connection.Close()
}

func (c *connector) closeChannel() error {
	if c.channel == nil {
		return nil
	}

	return c.channel.Close()
}

func (c *connector) connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.createConnection()
	if err != nil {
		return err
	}

	err = c.createChannel()
	if err != nil {
		return err
	}

	go c.reconnectJob(
		c.connection.NotifyClose(make(chan *amqp.Error, 1)),
		c.channel.NotifyClose(make(chan *amqp.Error, 1)),
	)

	return nil
}

func (c *connector) reconnect() {
	err := c.Close()
	if err != nil {
		c.logger.Err(err).Msg("failed to close previous connection")
	}

	for {
		err = c.connect()
		if err == nil {
			c.logger.Info().Msg("connection reconnected")
			return
		}
		c.logger.Err(err).Msg("failed to reconnect")
		time.Sleep(c.cfg.ReconnectDelay)
	}
}

func (c *connector) createConnection() error {
	conn, err := amqp.DialConfig(c.cfg.GetURI(), amqp.Config{
		Dial:      amqp.DefaultDial(c.cfg.ConnectTimeout),
		Heartbeat: c.cfg.Heartbeat,
		Locale:    defLocale,
	})
	if err != nil {
		return err
	}
	c.connection = conn

	return nil
}

func (c *connector) createChannel() error {
	channel, err := c.connection.Channel()
	if err != nil {
		return err
	}
	c.channel = channel

	return nil
}

func (c *connector) reconnectJob(connectionError, channelError <-chan *amqp.Error) {
	select {
	case err, ok := <-connectionError:
		if !ok {
			c.logger.Info().Msg("closing connector")
			return
		}
		if err != nil {
			c.logger.Err(err).Msg("received connection error")
		}
		c.reconnect()
	case err, ok := <-channelError:
		if !ok {
			c.logger.Info().Msg("closing connector")
			return
		}
		if err != nil {
			c.logger.Err(err).Msg("received channel error")
		}
		c.reconnect()
	}
}
