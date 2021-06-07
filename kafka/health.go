package kafka

import (
	"context"

	"github.com/city-mobil/gobuns/health"
)

func NewProducerHealthCheckCallback(p Producer) health.CheckCallback {
	return func(ctx context.Context) *health.CheckResult {
		status := health.CheckStatusFail
		if err := p.Ping(); err == nil {
			status = health.CheckStatusPass
		}

		return &health.CheckResult{
			ComponentID:   p.ComponentID(),
			ComponentType: p.ComponentType(),
			Status:        status,
		}
	}
}
