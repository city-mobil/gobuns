package health

import (
	"context"
	"time"
)

const (
	healthTimeCheckUnit        = "ms"
	healthCriticalResponseTime = 200 * time.Millisecond
)

// Checkable is an interface adapter for performing healthchecks.
type Checkable interface {
	// Ping is a method which is called for healthcheck.
	//
	// Ping does not have to 'ping' at all. It can do any check.
	Ping(context.Context) error

	// ComponentID is a id for some given Checkable.
	//
	// For example, it can be some 'uuid'.
	ComponentID() string

	// ComponentType is a type for some given Checkable.
	//
	// For example, it can be 'database'.
	ComponentType() string

	// Name is a name for some given Checkable.
	//
	// For example, it can be 'mysql'.
	Name() string
}

// NewResponseTimeCheckCallback creates new check callback.
//
// Option 'isSlave' defines if the given Checkable interface is a slave database connection.
// If Checkable is not a slave or master, isSlave must be set to false.
func NewResponseTimeCheckCallback(ch Checkable, isSlave bool) CheckCallback {
	return func(ctx context.Context) *CheckResult {
		res := &CheckResult{
			Status:        CheckStatusPass,
			ComponentID:   ch.ComponentID(),
			ComponentType: ch.ComponentType(),
			ObservedUnit:  healthTimeCheckUnit,
		}

		st := time.Now()
		err := ch.Ping(ctx)
		passed := time.Since(st)
		res.ObservedValue = passed.Milliseconds()

		// TODO(a.petrukhin): implement HealthCheckOptions().
		if passed > healthCriticalResponseTime {
			res.Status = CheckStatusWarn
		}

		if err == nil {
			return res
		}

		res.Error = err
		res.Output = err.Error()

		// NOTE(a.petrukhin): we consider that erroring slave is not a problem, because
		// there are many slaves for each master(usually).
		if isSlave {
			res.Status = CheckStatusWarn
		} else {
			res.Status = CheckStatusFail
		}
		return res
	}
}
