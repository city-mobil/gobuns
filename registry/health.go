package registry

import (
	"github.com/city-mobil/gobuns/health"
)

// NewHealthCheckCallback creates new health-check callback for local consul-agent
func NewHealthCheckCallback(client health.Checkable) health.CheckCallback {
	return health.NewResponseTimeCheckCallback(client, false)
}
