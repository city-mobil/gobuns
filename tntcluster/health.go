package tntcluster

import (
	"context"
	"fmt"

	"github.com/city-mobil/gobuns/health"
	"github.com/city-mobil/gobuns/tntcluster/pool"
)

type poolHealthAdapter struct {
	pool.ConnectorPool
}

func (p *poolHealthAdapter) Ping(_ context.Context) error {
	return p.ConnectorPool.Ping()
}

// NewClusterHealthCallbacks creates new health-check callbacks for a cluster.
func NewClusterHealthCallbacks(name string, c Cluster) (cbNames []string, cbs []health.CheckCallback, err error) {
	for _, shard := range c.GetShards() {
		names, callbacks, err := NewShardHealthCallbacks(name, shard)
		if err != nil {
			return nil, nil, err
		}
		cbNames = append(cbNames, names...)
		cbs = append(cbs, callbacks...)
	}

	return
}

// NewShardHealthCallbacks creates new health-check callbacks for a single shard.
func NewShardHealthCallbacks(name string, sh Shard) (cbNames []string, cbs []health.CheckCallback, err error) {
	var masterConn pool.ConnectorPool
	masterConn, err = sh.GetMasterConnector()
	if err != nil {
		err = fmt.Errorf("failed to get master connection for %s tarantool: %s", name, err)
		return
	}

	cbNames = append(cbNames, fmt.Sprintf("tarantool_%s:%s:master:responseTime", name, masterConn.RemoteAddr()))

	p := &poolHealthAdapter{
		ConnectorPool: masterConn,
	}
	cbs = append(cbs, health.NewResponseTimeCheckCallback(p, false))

	slaves, err := sh.GetSlaveConnectors()
	if err != nil {
		err = fmt.Errorf("failed to get replicas for %s tarantool: %s", name, err)
		return
	}
	for _, v := range slaves {
		cbNames = append(cbNames, fmt.Sprintf("tarantool_%s:%s:master:responseTime", name, v.RemoteAddr()))
		p := &poolHealthAdapter{
			ConnectorPool: v,
		}
		cbs = append(cbs, health.NewResponseTimeCheckCallback(p, true))
	}
	// NOTE(a.petrukhin): we track not average response time at all, but the ping response time.
	return
}
