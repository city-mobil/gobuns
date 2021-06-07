package promlib

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type InstrumentOption func(*instrument)

func InstrumentWithPath(fn func(r *http.Request) string) InstrumentOption {
	return func(i *instrument) {
		i.pathNameFunc = fn
	}
}

func InstrumentWithRegisterer(reg prometheus.Registerer) InstrumentOption {
	return func(i *instrument) {
		i.registry = reg
	}
}

type instrument struct {
	registry     prometheus.Registerer
	pathNameFunc func(r *http.Request) string
}

func newInstrument() *instrument {
	return &instrument{
		registry: prometheus.DefaultRegisterer,
		pathNameFunc: func(r *http.Request) string {
			return r.URL.Path
		},
	}
}

// InstrumentRoundTripper is a middleware that wraps the provided http.RoundTripper.
// It creates, registers and sets all needed metric handlers.
// Use this middleware to collect metrics from HTTP client.
//
// Middleware must be created only once for a given name.
func InstrumentRoundTripper(name string, next http.RoundTripper, opts ...InstrumentOption) http.RoundTripper {
	inst := newInstrument()
	for _, opt := range opts {
		opt(inst)
	}

	// An in-flight request is a request that has been started
	// but not yet completed, a.k.a. a request in progress.
	inFlightGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: name + "_in_flight_requests",
		Help: "A gauge of in-flight requests for HTTP client.",
	})

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name + "_requests_total",
			Help: "A counter for requests from the HTTP client.",
		},
		[]string{"code", "method", "path"},
	)

	// dnsLatencyVec uses custom buckets based on expected dns durations.
	// It has an instance label "event", which is set in the
	// DNSStart and DNSDone hook functions defined in the
	// InstrumentTrace struct below.
	dnsLatencyVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    name + "_dns_duration_seconds",
			Help:    "Trace DNS latency histogram.",
			Buckets: []float64{.05, .1, .25, .5},
		},
		[]string{"event"},
	)

	// tlsLatencyVec uses custom buckets based on expected tls durations.
	// It has an instance label "event", which is set in the
	// TLSHandshakeStart and TLSHandshakeDone hook functions defined in the
	// InstrumentTrace struct below.
	tlsLatencyVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    name + "_tls_duration_seconds",
			Help:    "Trace TLS latency histogram.",
			Buckets: []float64{.05, .1, .25, .5},
		},
		[]string{"event"},
	)

	histVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    name + "_request_duration_seconds",
			Help:    "A histogram of request latencies.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	inst.registry.MustRegister(counter, tlsLatencyVec, dnsLatencyVec, histVec, inFlightGauge)

	// Define functions for the available httptrace.ClientTrace hook
	// functions that we want to instrument.
	trace := &promhttp.InstrumentTrace{
		DNSStart: func(t float64) {
			dnsLatencyVec.WithLabelValues("dns_start").Observe(t)
		},
		DNSDone: func(t float64) {
			dnsLatencyVec.WithLabelValues("dns_done").Observe(t)
		},
		TLSHandshakeStart: func(t float64) {
			tlsLatencyVec.WithLabelValues("tls_handshake_start").Observe(t)
		},
		TLSHandshakeDone: func(t float64) {
			tlsLatencyVec.WithLabelValues("tls_handshake_done").Observe(t)
		},
	}

	return promhttp.InstrumentRoundTripperInFlight(inFlightGauge,
		inst.instrumentRoundTripperCounter(counter,
			promhttp.InstrumentRoundTripperTrace(trace,
				inst.instrumentRoundTripperDuration(histVec, next),
			),
		),
	)
}

func (i *instrument) instrumentRoundTripperCounter(counter *prometheus.CounterVec, next http.RoundTripper) promhttp.RoundTripperFunc {
	return func(r *http.Request) (*http.Response, error) {
		resp, err := next.RoundTrip(r)
		if err == nil {
			labels := prometheus.Labels{
				"code":   sanitizeCode(resp.StatusCode),
				"method": sanitizeMethod(r.Method),
				"path":   i.pathNameFunc(r),
			}

			counter.With(labels).Inc()
		}
		return resp, err
	}
}

func (i *instrument) instrumentRoundTripperDuration(obs prometheus.ObserverVec, next http.RoundTripper) promhttp.RoundTripperFunc {
	return func(r *http.Request) (*http.Response, error) {
		start := time.Now()
		resp, err := next.RoundTrip(r)
		if err == nil {
			labels := prometheus.Labels{
				"method": sanitizeMethod(r.Method),
				"path":   i.pathNameFunc(r),
			}

			obs.With(labels).Observe(time.Since(start).Seconds())
		}
		return resp, err
	}
}

func sanitizeMethod(m string) string {
	switch m {
	case "GET", "get":
		return "get"
	case "PUT", "put":
		return "put"
	case "HEAD", "head":
		return "head"
	case "POST", "post":
		return "post"
	case "DELETE", "delete":
		return "delete"
	case "CONNECT", "connect":
		return "connect"
	case "OPTIONS", "options":
		return "options"
	case "NOTIFY", "notify":
		return "notify"
	default:
		return strings.ToLower(m)
	}
}

// If the wrapped http.Handler has not set a status code, i.e. the value is
// currently 0, sanitizeCode will return 200, for consistency with behavior in the stdlib.
func sanitizeCode(s int) string {
	switch s {
	case 100:
		return "100"
	case 101:
		return "101"

	case 200, 0:
		return "200"
	case 201:
		return "201"
	case 202:
		return "202"
	case 203:
		return "203"
	case 204:
		return "204"
	case 205:
		return "205"
	case 206:
		return "206"

	case 300:
		return "300"
	case 301:
		return "301"
	case 302:
		return "302"
	case 304:
		return "304"
	case 305:
		return "305"
	case 307:
		return "307"

	case 400:
		return "400"
	case 401:
		return "401"
	case 402:
		return "402"
	case 403:
		return "403"
	case 404:
		return "404"
	case 405:
		return "405"
	case 406:
		return "406"
	case 407:
		return "407"
	case 408:
		return "408"
	case 409:
		return "409"
	case 410:
		return "410"
	case 411:
		return "411"
	case 412:
		return "412"
	case 413:
		return "413"
	case 414:
		return "414"
	case 415:
		return "415"
	case 416:
		return "416"
	case 417:
		return "417"
	case 418:
		return "418"

	case 500:
		return "500"
	case 501:
		return "501"
	case 502:
		return "502"
	case 503:
		return "503"
	case 504:
		return "504"
	case 505:
		return "505"

	case 428:
		return "428"
	case 429:
		return "429"
	case 431:
		return "431"
	case 511:
		return "511"

	default:
		return strconv.Itoa(s)
	}
}
