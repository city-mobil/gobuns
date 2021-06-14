package kafka

import (
	"context"

	"github.com/city-mobil/gobuns/health"
)

func NewProducerHealthCheckCallback(p Producer) health.CheckCallback {
	return func(ctx context.Context) *health.CheckResult {
		status := health.CheckStatusFail
		err := p.Ping()
		if err == nil {
			status = health.CheckStatusPass
		}

		return &health.CheckResult{
			ComponentID:   p.ComponentID(),
			ComponentType: p.ComponentType(),
			Status:        status,
			Error:         err,
		}
	}
}

func NewConsumerHealthCheckCallback(c Consumer) health.CheckCallback {
	return func(ctx context.Context) *health.CheckResult {
		status := health.CheckStatusFail
		err := c.Ping()
		if err == nil {
			status = health.CheckStatusPass
		}

		return &health.CheckResult{
			ComponentID:   c.ComponentID(),
			ComponentType: c.ComponentType(),
			Status:        status,
			Error:         err,
		}
	}
}
