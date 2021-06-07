// Package external contains http client for some external HTTP requests.
package external

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"

	"github.com/city-mobil/gobuns/promlib"
	"github.com/city-mobil/gobuns/retry"
)

const (
	traceComponentName = "go-buns/external"
)

var (
	ErrEmptyName = errors.New("client name must be set to collect metrics")
)

// Client is an external HTTP client interface.
type Client interface {
	// SetCustomTransport sets custom HTTP transport for a client.
	//
	// Be aware that custom transport overrides
	// all open tracing and prometheus middlewares.
	SetCustomTransport(http.RoundTripper)

	// Get performs GET-requests with retries for prepared http.Request
	Get(context.Context, *http.Request) (*http.Response, error)

	// Post performs POST-requests with retries for prepared http.Request
	Post(context.Context, *http.Request) (*http.Response, error)

	// Do performs one given http.Request with retries.
	Do(context.Context, *http.Request) (*http.Response, error)
}

type client struct {
	client  *http.Client
	retrier *retry.Retrier
	onRetry func(n uint, err error)
}

var (
	// DefaultClient is an external.Client initialized with default configuration and used for Get, Post and Do.
	DefaultClient Client
)

func init() {
	cfg := &Config{}

	var err error
	DefaultClient, err = New(cfg)
	if err != nil {
		panic(err)
	}
}

// New creates a new Client for a given Config.
func New(userCfg *Config) (Client, error) {
	cfg := userCfg.withDefaults()

	if cfg.Metrics.Collect && cfg.Name == "" {
		return nil, ErrEmptyName
	}

	tr, err := newDefaultTransport(&cfg)
	if err != nil {
		return nil, err
	}

	if cfg.Metrics.Collect {
		tr = promlib.InstrumentRoundTripper(cfg.Name, tr, cfg.Metrics.Options...)
	}

	return &client{
		client: &http.Client{
			Transport: &nethttp.Transport{
				RoundTripper: tr,
			},
			Timeout: cfg.RequestTimeout,
		},
		retrier: retry.New(cfg.RetryConfig),
		onRetry: cfg.OnRetry,
	}, nil
}

func (c *client) SetCustomTransport(tr http.RoundTripper) {
	c.client.Transport = tr
}

// Get performs GET-requests with retries for prepared http.Request
//
// Get also checks if the given request is a real GET request, otherwise an error is returned.
func (c *client) Get(ctx context.Context, r *http.Request) (*http.Response, error) {
	if r.Method != http.MethodGet {
		return nil, fmt.Errorf("Invalid request method specified: %s, expected GET", r.Method) //nolint:golint,stylecheck
	}
	return c.doRequest(ctx, r)
}

// Post performs POST-requests with retries for prepared http.Request
//
// Post also checks if the given request is a real POST request, otherwise an error is returned.
func (c *client) Post(ctx context.Context, r *http.Request) (*http.Response, error) {
	if r.Method != http.MethodPost {
		return nil, fmt.Errorf("Invalid request method specified: %s, expected POST", r.Method) //nolint:golint,stylecheck
	}
	return c.doRequest(ctx, r)
}

// Do performs one given http.Request with retries.
func (c *client) Do(ctx context.Context, r *http.Request) (*http.Response, error) {
	return c.doRequest(ctx, r)
}

func (c *client) doRequest(ctx context.Context, r *http.Request) (resp *http.Response, err error) {
	var (
		ht   *nethttp.Tracer
		body io.ReadCloser
	)

	r = r.WithContext(httptrace.WithClientTrace(ctx, &httptrace.ClientTrace{}))
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		r, ht = nethttp.TraceRequest(span.Tracer(), r, nethttp.ComponentName(traceComponentName))
		defer ht.Finish()
	}

	action := func() error {
		// NOTE(a.petrukhin): here we perform copy of current body.
		// It is faster and consumes less memory to give raw bytes as body and create new reader here.
		if r.Body != nil && r.GetBody != nil {
			body, err = r.GetBody()
			if err != nil {
				return retry.Unrecoverable(fmt.Errorf("get request body error: %s", err))
			}
		}
		// If body is provided, we just recreate body in order to have possibility to retry.
		// For further information see https://stackoverflow.com/questions/31337891/net-http-http-contentlength-222-with-body-length-0
		//
		// We manually close the body on retry, otherwise the client must do it himself.
		resp, err = c.client.Do(r) //nolint:bodyclose
		if err == nil {
			return nil
		}

		// NOTE(a.petrukhin): setting back because the old body was probably half-read.
		r.Body = body

		isRetryable := false
		if v, ok := err.(net.Error); ok && v.Timeout() {
			// NOTE(a.petrukhin): re-attempt if request timed out.
			isRetryable = true
		}

		if isRetryable {
			return err
		}

		return retry.Unrecoverable(err)
	}

	onRetry := func(n uint, err error) {
		if c.onRetry != nil {
			c.onRetry(n, err)
		}

		if resp != nil {
			// If all attempts were failed,
			// do not close the body for the last response.
			// It is the user responsibility.
			_ = resp.Body.Close()
		}
	}

	err = c.retrier.Do(ctx, action, onRetry)
	if err != nil && span != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.Error(err))
	}

	return resp, err
}

// Get performs GET-requests with retries for prepared http.Request
//
// Get also checks if the given request is a real GET request, otherwise an error is returned.
func Get(ctx context.Context, r *http.Request) (*http.Response, error) {
	return DefaultClient.Get(ctx, r)
}

// Post performs POST-requests with retries for prepared http.Request
//
// Post also checks if the given request is a real POST request, otherwise an error is returned.
func Post(ctx context.Context, r *http.Request) (*http.Response, error) {
	return DefaultClient.Post(ctx, r)
}

// Do performs one given http.Request with retries.
func Do(ctx context.Context, r *http.Request) (*http.Response, error) {
	return DefaultClient.Do(ctx, r)
}
