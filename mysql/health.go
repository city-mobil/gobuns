package mysql

import (
	"context"
	"fmt"

	"github.com/city-mobil/gobuns/health"
)

type healthAdapter struct {
	Adapter
}

func (h *healthAdapter) Ping(ctx context.Context) error {
	return h.PingContext(ctx)
}

// NewHealthCheckCallbacks creates new callbacks for health-check.
func NewHealthCheckCallbacks(sh Shard) (names []string, cbs []health.CheckCallback) {
	// NOTE(a.petrukhin): no #BlackLivesMatter here!
	names = append(names, fmt.Sprintf("mysql:%s:master:responseTime", sh.GetMasterConn().Name()))
	// NOTE(a.petrukhin): we track not average response time at all, but the ping response time.
	cbs = append(cbs, health.NewResponseTimeCheckCallback(&healthAdapter{
		Adapter: sh.GetMasterConn(),
	}, false))
	for i, slave := range sh.GetSlaveConnections() {
		names = append(names, fmt.Sprintf("mysql_%s:%s:slave_%d:responseTime", slave.Name(), slave.ComponentID(), i))
		cbs = append(cbs, health.NewResponseTimeCheckCallback(&healthAdapter{
			Adapter: slave,
		}, true))
	}
	return
}
