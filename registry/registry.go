package registry

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/city-mobil/gobuns/retry"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

var (
	ErrKeyNotExist = errors.New("key does not exist")
)

const (
	defaultName        = "agent"
	defaultType        = "consul"
	traceComponentName = "go-buns/registry"
)

type Client interface {
	GetString(context.Context, string) (string, error)
	GetBool(context.Context, string) (bool, error)
	GetInt(context.Context, string) (int, error)

	Ping(context.Context) error
	Name() string
	ComponentType() string
	ComponentID() string

	SetName(string)
	SetComponentType(string)
	SetComponentID(string)
}

type client struct {
	consulClient *consulapi.Client
	queryTimeout time.Duration
	retrier      *retry.Retrier

	name          string
	componentType string
	componentID   string
}

func NewClient(config Config) (Client, error) {
	consulCfg := consulapi.DefaultConfig()
	consulCfg.Address = config.Addr

	c, err := consulapi.NewClient(consulCfg)
	if err != nil {
		return nil, err
	}

	return &client{
		consulClient: c,
		queryTimeout: config.QueryTimeout,
		retrier:      retry.New(config.RetryConfig),

		name:          defaultName,
		componentType: defaultType,
		componentID:   config.Addr,
	}, nil
}

func (c *client) GetString(ctx context.Context, key string) (string, error) {
	val, err := c.getWithRetries(ctx, key)
	if err != nil {
		return "", err
	}

	return val, err
}

func (c *client) GetBool(ctx context.Context, key string) (bool, error) {
	val, err := c.getWithRetries(ctx, key)
	if err != nil {
		return false, err
	}

	return strconv.ParseBool(val)
}

func (c *client) GetInt(ctx context.Context, key string) (int, error) {
	val, err := c.getWithRetries(ctx, key)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(val)
}

func (c *client) Ping(_ context.Context) error {
	hc, _, err := c.consulClient.Health().Node(c.componentID, nil)
	if err != nil {
		return err
	}

	status := hc.AggregatedStatus()
	if status != consulapi.HealthPassing {
		return fmt.Errorf("consulagent unexpected healthcheck status: %s", status)
	}

	return nil
}

func (c *client) Name() string {
	return c.name
}

func (c *client) ComponentType() string {
	return c.componentType
}

func (c *client) ComponentID() string {
	return c.componentID
}

func (c *client) SetName(name string) {
	c.name = name
}

func (c *client) SetComponentType(componentType string) {
	c.componentType = componentType
}

func (c *client) SetComponentID(componentID string) {
	c.componentID = componentID
}

func (c *client) getWithRetries(ctx context.Context, key string) (string, error) {
	var (
		resp *consulapi.KVPair
		span opentracing.Span
		err  error
	)
	rootSpan := opentracing.SpanFromContext(ctx)

	action := func() error {
		if rootSpan != nil {
			if span == nil {
				span, _ = opentracing.StartSpanFromContextWithTracer(ctx, rootSpan.Tracer(), "consul/get")
			} else {
				span = rootSpan.Tracer().StartSpan("consul/get", opentracing.FollowsFrom(span.Context()))
			}
			span.SetTag("key", key)
			ext.Component.Set(span, traceComponentName)
			ext.SpanKindRPCClient.Set(span)
			ext.PeerHostname.Set(span, c.componentID)
		}
		resp, err = c.get(ctx, key)
		if span != nil {
			if err != nil {
				ext.Error.Set(span, true)
				span.LogFields(log.Error(err))
			}
			span.Finish()
		}
		if err != nil {
			if consulapi.IsRetryableError(err) {
				return err
			}

			return retry.Unrecoverable(err)
		}

		return nil
	}

	onRetry := func(n uint, err error) {
		// nothing to do
	}

	err = c.retrier.Do(ctx, action, onRetry)
	if err != nil {
		return "", err
	}

	if resp == nil {
		return "", ErrKeyNotExist
	}

	return string(resp.Value), nil
}

func (c *client) get(ctx context.Context, key string) (*consulapi.KVPair, error) {
	execCtx := ctx
	var cancel context.CancelFunc = func() {}

	if c.queryTimeout > 0 {
		execCtx, cancel = context.WithTimeout(ctx, c.queryTimeout)
	}
	defer cancel()

	queryOptions := &consulapi.QueryOptions{}
	queryOptions = queryOptions.WithContext(execCtx)

	resp, _, err := c.consulClient.KV().Get(key, queryOptions)

	return resp, err
}
