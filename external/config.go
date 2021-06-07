package external

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"time"

	"github.com/city-mobil/gobuns/config"
	"github.com/city-mobil/gobuns/promlib"
	"github.com/city-mobil/gobuns/retry"
)

type VersionTLS string

const (
	VersionTLS10 VersionTLS = "1.0"
	VersionTLS11 VersionTLS = "1.1"
	VersionTLS12 VersionTLS = "1.2"
	VersionTLS13 VersionTLS = "1.3"
)

const (
	defDialTimeout       = 5 * time.Second
	defIdleTimeout       = 15 * time.Minute
	defRequestTimeout    = 500 * time.Millisecond
	defKeepAliveInterval = 30 * time.Second
	defMaxIdleConns      = 100
	defNoHTTPS           = true
	defVersionTLS        = VersionTLS12
)

var (
	defOnRetry = func(n uint, err error) {
		log.Printf("[external] request error: %s, attempt: %d", err, n)
	}
)

// Config is a external HTTP client configuration.
type Config struct {
	// Name is a unique name of the client.
	Name string

	// DialContext connects to the address on the named network using
	// the provided context.
	DialContext func(ctx context.Context, network, addr string) (net.Conn, error)

	// MaxIdleConnsPerHost, if non-zero, controls the maximum idle
	// (keep-alive) connections to keep per-host. If zero,
	// defMaxIdleConns is used.
	MaxIdleConns int

	// IdleConnTimeout is the maximum amount of time an idle
	// (keep-alive) connection will remain idle before closing
	// itself.
	// Zero means no limit.
	IdleConnTimeout time.Duration

	// RequestTimeout specifies a time limit for requests made by this
	// Client. The timeout includes connection time, any
	// redirects, and reading the response body. The timer remains
	// running after Get, Head, Post, or Do return and will
	// interrupt reading of the Response.Body.
	//
	// A Timeout of zero means no timeout.
	//
	// The Client cancels requests to the underlying Transport
	// as if the Request's Context ended.
	//
	// For compatibility, the Client will also use the deprecated
	// CancelRequest method on Transport if found. New
	// RoundTripper implementations should use the Request's Context
	// for cancellation instead of implementing CancelRequest.
	RequestTimeout time.Duration

	// RetryConfig is a configuration for retry policy.
	RetryConfig *retry.Config

	OnRetry func(n uint, err error)

	// MinVersionTLS contains the minimum TLS version that is acceptable.
	MinVersionTLS VersionTLS

	// PublicCertPath is a OS-path for the public HTTPS certificate.
	PublicCertPath string

	// PublicCertPath is a OS-path for the private HTTPS certificate.
	PrivateCertPath string

	// NoHTTPS ignores HTTPS if set.
	NoHTTPS bool

	// ForceInsecureSkipVerify forces to set 'InsecureSkipVerify' even if
	// all certificates are set.
	//
	// This option needs to be set when error 'x509: certificate signed by unknown authority'
	// is occurred.
	ForceInsecureSkipVerify bool

	// Metrics defines options how to collect metrics in Prometheus format.
	Metrics ConfigMetrics
}

type ConfigMetrics struct {
	// Collect enables gathering metrics
	// of the HTTP client in Prometheus format.
	//
	// Name of the client must be set and unique.
	Collect bool

	// Options allows to customize Prometheus metric collector.
	Options []promlib.InstrumentOption
}

// NewConfig is a new config callback with given prefix.
// All the config variables MUST be registered before the callback is called.
//
// It can be used in such way:
//  // cfg := NewConfig("some_prefix")
//  // config.InitOnce()
//  // client := New(cfg())
//
func NewConfig(prefix string) func() *Config {
	if prefix != "" {
		prefix += ".external."
	} else {
		prefix = "external."
	}

	p := func(opt string) string {
		return prefix + opt
	}

	var (
		clientName              = config.String(p("client_name"), "", "unique name of the client")
		dialTimeout             = config.Duration(p("dial_timeout"), defDialTimeout, "dial timeout")
		keepAlive               = config.Duration(p("keepalive_interval"), defKeepAliveInterval, "keepalive messages interval (protocol specific)")
		maxIdleConns            = config.Int(p("max_idle_conns"), defMaxIdleConns, "max idle connections to external service")
		idleTimeout             = config.Duration(p("idle_conn_timeout"), defIdleTimeout, "idle connection timeout")
		requestTimeout          = config.Duration(p("request_timeout"), defRequestTimeout, "request timeout")
		noHTTPS                 = config.Bool(p("no_https"), defNoHTTPS, "controls whether a client verifies the server's certificate chain and host name")
		forceInsecureSkipVerify = config.Bool(p("force_insecure_skip_verify"), false, "force insecure skip verify option")
		retryCfgFn              = retry.GetRetryConfig(p("retries"))
		tlsVersion              = config.String(p("tls.version"), "1.2", "TLS version")
		tlsPublicCert           = config.String(p("tls.cert.public"), "", "path to a public client TLS cert")
		tlsPrivateCert          = config.String(p("tls.cert.private"), "", "path to a private client TLS cert")
		metricsCollect          = config.Bool(p("metrics.collect"), false, "enables gathering metrics in Prometheus format")
	)

	return func() *Config {
		return &Config{
			Name:                    *clientName,
			DialContext:             makeDialContext(*dialTimeout, *keepAlive),
			MaxIdleConns:            *maxIdleConns,
			IdleConnTimeout:         *idleTimeout,
			RequestTimeout:          *requestTimeout,
			NoHTTPS:                 *noHTTPS,
			RetryConfig:             retryCfgFn(),
			MinVersionTLS:           VersionTLS(*tlsVersion),
			PrivateCertPath:         *tlsPrivateCert,
			ForceInsecureSkipVerify: *forceInsecureSkipVerify,
			PublicCertPath:          *tlsPublicCert,
			Metrics: ConfigMetrics{
				Collect: *metricsCollect,
			},
		}
	}
}

// withDefaults sets default parameters for config if some are not set.
//
// If the Config is nil, new Config is created and filled with default params.
func (cfg *Config) withDefaults() (c Config) {
	if cfg != nil {
		c = *cfg
	}

	if c.DialContext == nil {
		c.DialContext = (new(net.Dialer)).DialContext
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = defMaxIdleConns
	}
	if c.RetryConfig == nil {
		c.RetryConfig = retry.NewDefRetryConfig()
	}
	if c.RequestTimeout == 0 {
		c.RequestTimeout = defRequestTimeout
	}
	if c.OnRetry == nil {
		c.OnRetry = defOnRetry
	}
	if c.MinVersionTLS == "" {
		c.MinVersionTLS = defVersionTLS
	}

	return
}

func makeDialContext(timeout, keepalive time.Duration) func(ctx context.Context, network, addr string) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   timeout,
		KeepAlive: keepalive,
	}

	return dialer.DialContext
}

func castTLSVersion(v VersionTLS) uint16 {
	switch v {
	case VersionTLS10:
		return tls.VersionTLS10
	case VersionTLS11:
		return tls.VersionTLS11
	case VersionTLS12:
		return tls.VersionTLS12
	case VersionTLS13:
		return tls.VersionTLS13
	default:
		return castTLSVersion(defVersionTLS)
	}
}
