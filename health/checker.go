package health

import (
	"context"
	"fmt"
	"sync"
)

// CheckCallback is a callback with is called during each healthcheck.
type CheckCallback func(ctx context.Context) *CheckResult

// Checker is an interface which provides healthcheck methods.
type Checker interface {
	// AddCallback adds a single callback in healthcheck list.
	AddCallback(string, CheckCallback)

	// AddMultipleCallbacks adds multiple callbacks in healthcheck list.
	AddMultipleCallbacks([]string, []CheckCallback)

	// Check performs a single healthcheck.
	Check() *CheckResponse

	CheckContext(context.Context) *CheckResponse
}

// CheckerOptions defines options for a single checker.
type CheckerOptions struct {
	// Version is a version of checker.
	//
	// For example, it can be commit hash.
	Version string

	// ReleaseID is a release id of checker.
	//
	// For example, it can be git tag of current release.
	ReleaseID string

	// ServiceID is a service id of checker.
	//
	// For example, it can be ip address of current host or current hostname.
	ServiceID string
}

type checker struct {
	mu        sync.RWMutex
	callbacks []callback
	opts      CheckerOptions
}

// NewChecker creates new Checker.
func NewChecker(opts CheckerOptions) Checker {
	return &checker{
		opts: opts,
	}
}

type callback struct {
	name string
	cb   CheckCallback
}

// AddCallback adds single callback function for healthcheck.
func (c *checker) AddCallback(name string, cb CheckCallback) {
	c.mu.Lock()
	c.callbacks = append(c.callbacks, callback{
		name: name,
		cb:   cb,
	})
	c.mu.Unlock()
}

// AddMultipleCallbacks adds multiple callback function for healthcheck.
func (c *checker) AddMultipleCallbacks(names []string, cbs []CheckCallback) {
	if len(names) != len(cbs) {
		panic(fmt.Sprintf("names length and callbacks length are not equal: %d != %d", len(names), len(cbs)))
	}
	for i := 0; i < len(names); i++ {
		c.AddCallback(names[i], cbs[i])
	}
}

// CheckContext performs a single healthcheck for previously added callbacks with given context.
func (c *checker) CheckContext(ctx context.Context) *CheckResponse {
	return c.check(ctx)
}

// Check performs a single healthcheck for previously added callbacks.
func (c *checker) Check() *CheckResponse {
	return c.check(context.Background())
}

func (c *checker) check(ctx context.Context) *CheckResponse {
	// NOTE(a.petrukhin): it is a race. But it is a by-design race.
	// We can not delete callbacks, we can only add them. If one is added after the
	// lock is taken, the callback is not going to be called during the current check, but is going to be called during the next one.
	c.mu.RLock()
	l := len(c.callbacks)
	c.mu.RUnlock()
	callbacks := make([]callback, l)
	copy(callbacks, c.callbacks)

	result := &CheckResponse{
		ReleaseID: c.opts.ReleaseID,
		ServiceID: c.opts.ServiceID,
		Version:   c.opts.Version,
		Checks:    make(map[string]*CheckResult, l),
		Status:    CheckStatusPass,
	}
	var wg sync.WaitGroup
	var mu sync.RWMutex
	wg.Add(l)
	for _, cb := range callbacks {
		go func(cb callback) {
			defer wg.Done()
			res := handleCallback(ctx, cb.cb)
			if res == nil {
				return // NOTE: if callback return nil, we would not receive any report about specific check
			}
			mu.Lock()
			if res.Status == CheckStatusFail {
				result.Status = CheckStatusFail
				result.Output = res.Output
			}
			result.Checks[cb.name] = res
			mu.Unlock()
		}(cb)
	}

	wg.Wait()
	return result
}

func handleCallback(ctx context.Context, cb CheckCallback) (res *CheckResult) {
	defer func() {
		if r := recover(); r != nil {
			res = &CheckResult{
				Status: CheckStatusFail,
				Output: fmt.Sprintf("%v", r),
			}
		}
	}()

	res = cb(ctx)
	if res == nil {
		return
	}

	if res.Error != nil {
		res.Output = res.Error.Error()
		switch res.Error.(type) {
		case *FailError:
			res.Status = CheckStatusFail
		case *WarnError:
			res.Status = CheckStatusWarn
		}
	}

	return
}
