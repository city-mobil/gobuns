package handlers

import (
	"net/http"
	"time"

	"github.com/city-mobil/gobuns/httputil"
	"github.com/city-mobil/gobuns/zlog"
	"github.com/city-mobil/gobuns/zlog/hlog"
)

// Filter is a filter which decides when access logger
// should log an incoming request or not.
//
// Returns true if a request will be logged, otherwise false.
type Filter func(code int, dur time.Duration, err error) bool

// LogAll is a filter to log any requests.
func LogAll(_ int, _ time.Duration, _ error) bool {
	return true
}

// Log5xx is a filter to log only requests with response status code 5xx.
func Log5xx(code int, _ time.Duration, _ error) bool {
	return code >= http.StatusInternalServerError && code <= 599
}

// LogExcept2xx is a filter to log all requests with non-2xx response status code.
func LogExcept2xx(code int, _ time.Duration, _ error) bool {
	return code < http.StatusOK || code >= http.StatusMultipleChoices
}

type Option func(aL *AccessLogger)

// WithLogger sets provided logger to use by access logger.
func WithLogger(logger zlog.Logger) Option {
	return func(aL *AccessLogger) {
		aL.logger = logger
	}
}

// WithLoggerFromReq specifies to extract
// the logger instance from the request context.
//
// This option is not compatible with WithLogger option.
func WithLoggerFromReq() Option {
	return func(aL *AccessLogger) {
		aL.loggerFromReq = true
	}
}

// WithFilter sets the incoming log filter.
func WithFilter(f Filter) Option {
	return func(aL *AccessLogger) {
		aL.filter = f
	}
}

// WithIPLookup sets the custom IP lookup.
func WithIPLookup(ipl *httputil.IPLookup) Option {
	return func(aL *AccessLogger) {
		aL.ipLookup = ipl
	}
}

type AccessLogger struct {
	logger        zlog.Logger
	endpoint      string
	filter        Filter
	ipLookup      *httputil.IPLookup
	loggerFromReq bool
}

// NewAccessLogger creates new access logger with provided options.
func NewAccessLogger(endpoint string, opts ...Option) *AccessLogger {
	aL := &AccessLogger{
		endpoint: endpoint,
	}

	for _, opt := range opts {
		opt(aL)
	}

	if aL.logger == nil {
		aL.logger = zlog.Nop()
	}

	if aL.filter == nil {
		aL.filter = LogAll
	}

	if aL.ipLookup == nil {
		aL.ipLookup = httputil.NewIPLookup()
	}

	return aL
}

type response struct {
	dur  time.Duration
	code int
	err  error
}

func (r *response) hasErr() bool {
	return r.code >= 500 || r.err != nil
}

func (aL *AccessLogger) logRequest(req *http.Request, resp *response) {
	if !aL.filter(resp.code, resp.dur, resp.err) {
		return
	}

	logLevel := zlog.InfoLevel
	if resp.hasErr() {
		logLevel = zlog.ErrorLevel
	}

	ip := aL.ipLookup.GetRemoteIP(req)

	ev := aL.lg(req).WithLevel(logLevel).
		Dur("request_duration", resp.dur).
		Int("http_status", resp.code).
		Str("user_ip", ip).
		Str("query", req.Method+" "+req.URL.String()).
		AnErr("response_error", resp.err)

	if !aL.loggerFromReq {
		ev.Str("method", req.Method)
	}

	ev.Send()
}

func (aL *AccessLogger) error(r *http.Request, msg string, err error) {
	aL.lg(r).Err(err).Msg(msg)
}

func (aL *AccessLogger) warn(r *http.Request, msg string, err error) {
	aL.lg(r).Warn().Err(err).Msg(msg)
}

func (aL *AccessLogger) lg(r *http.Request) zlog.Logger {
	logger := aL.logger
	if aL.loggerFromReq {
		logger = hlog.FromRequest(r)
	}

	return logger
}
