package promlib

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	DefHTTPRequestDurBuckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
)

type Option func(*httpMiddleware)

func WithCustomPath(fn func(r *http.Request) string) Option {
	return func(middleware *httpMiddleware) {
		middleware.pathNameFunc = fn
	}
}

func WithUserAgentLabel() Option {
	return func(middleware *httpMiddleware) {
		middleware.useUserAgent = true
	}
}

func WithHistogramName(name string) Option {
	return func(middleware *httpMiddleware) {
		middleware.name = name
	}
}

type HTTPMiddleware interface {
	Handler(http.Handler) http.Handler
	HandlerFunc(http.HandlerFunc) http.HandlerFunc
}

type httpMiddleware struct {
	requestDurHistogram *prometheus.HistogramVec
	pathNameFunc        func(r *http.Request) string
	useUserAgent        bool
	name                string
}

func NewMiddleware(durationBuckets []float64, opts ...Option) HTTPMiddleware {
	if len(durationBuckets) == 0 {
		durationBuckets = DefHTTPRequestDurBuckets
	}

	middleware := newDefaultMW()

	for _, opt := range opts {
		opt(middleware)
	}

	labels := []string{"path", "code", "method"}
	if middleware.useUserAgent {
		labels = append(labels, "agent")
	}

	httpRequestDurHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    middleware.name,
		Help:    "Duration of HTTP request",
		Buckets: durationBuckets,
	}, labels)
	prometheus.MustRegister(httpRequestDurHistogram)

	middleware.requestDurHistogram = httpRequestDurHistogram
	return middleware
}

func newDefaultMW() *httpMiddleware {
	return &httpMiddleware{
		pathNameFunc: func(r *http.Request) string {
			return r.URL.Path
		},
		useUserAgent: false,
		name:         "http_request_duration_seconds",
	}
}

func (m *httpMiddleware) Handler(next http.Handler) http.Handler {
	return m.HandlerFunc(next.ServeHTTP)
}

func (m *httpMiddleware) HandlerFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		labels := make(prometheus.Labels, 2)
		labels["path"] = m.pathNameFunc(r)
		if m.useUserAgent {
			labels["agent"] = r.UserAgent()
		}

		handler := promhttp.InstrumentHandlerDuration(
			m.requestDurHistogram.MustCurryWith(labels),
			next,
		)

		handler(w, r)
	}
}
