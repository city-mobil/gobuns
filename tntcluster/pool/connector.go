package pool

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/city-mobil/gobuns/tntcluster/once"
	"github.com/viciious/go-tarantool"
)

type ConnectorPool interface {
	Close()
	Connect() (*tarantool.Connection, error)
	RemoteAddr() string

	Ping() error
	Name() string
	ComponentType() string
	ComponentID() string
}

const (
	componentTypeTarantool = "tarantool"
)

type connectorPool struct {
	connectors        []*tarantool.Connector
	remoteAddr        string
	name              string
	options           *tarantool.Options
	poolSize          int
	lastUsedConnector uint32
	once              once.Once
}

func New(dsnString, name string, options *tarantool.Options, poolSize int) ConnectorPool {
	return &connectorPool{
		name:       name,
		remoteAddr: dsnString,
		options:    options,
		poolSize:   poolSize,
		connectors: make([]*tarantool.Connector, 0, poolSize),
	}
}

func (c *connectorPool) Connect() (*tarantool.Connection, error) {
	connector := c.getNextConnector()
	return connector.Connect()
}

func (c *connectorPool) Ping() error {
	conn := c.getNextConnector()
	cn, err := conn.Connect()
	if err != nil {
		return err
	}

	res := cn.Exec(context.Background(), &tarantool.Ping{})
	if res == nil {
		return errors.New("got nil response")
	}
	return res.Error
}

func (c *connectorPool) ComponentID() string {
	return c.remoteAddr
}

func (c *connectorPool) Name() string {
	return c.name
}

func (c *connectorPool) ComponentType() string {
	return componentTypeTarantool
}

func (c *connectorPool) initPool() {
	var connector *tarantool.Connector
	connectors := make([]*tarantool.Connector, 0, c.poolSize)
	for i := 0; i < c.poolSize; i++ {
		connector = tarantool.New(c.remoteAddr, c.options)
		connectors = append(connectors, connector)
	}
	c.connectors = connectors
}

func (c *connectorPool) getNextConnector() *tarantool.Connector {
	c.once.Do(c.initPool)
	return c.getRoundRobinConnector()
}

func (c *connectorPool) RemoteAddr() string {
	return c.remoteAddr
}

func (c *connectorPool) getRoundRobinConnector() *tarantool.Connector {
	if len(c.connectors) == 1 {
		return c.connectors[0]
	}
	next := atomic.AddUint32(&c.lastUsedConnector, 1)
	idx := (int(next) - 1) % len(c.connectors)
	return c.connectors[idx]
}

func (c *connectorPool) Close() {
	connectors := c.connectors
	c.once.Reset()

	for _, connector := range connectors {
		connector.Close()
	}
}
