package tracing

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type HTTMMiddleware interface {
	Handler(h http.Handler) http.Handler
	HandlerFunc(h http.HandlerFunc) http.HandlerFunc
}

type Option func(*httpMiddleware)

func OperationName(fn func(r *http.Request) string) Option {
	return func(middleware *httpMiddleware) {
		middleware.opNameFunc = fn
	}
}

type httpMiddleware struct {
	tracer     opentracing.Tracer
	opNameFunc func(r *http.Request) string
}

func NewHTTPMiddleware(tracer opentracing.Tracer, opts ...Option) HTTMMiddleware {
	middleware := &httpMiddleware{
		tracer: tracer,
		opNameFunc: func(r *http.Request) string {
			return "HTTP " + r.Method + ":" + r.URL.Path
		},
	}
	for _, opt := range opts {
		opt(middleware)
	}
	return middleware
}

func (wr *httpMiddleware) Handler(h http.Handler) http.Handler {
	return wr.HandlerFunc(h.ServeHTTP)
}

func (wr *httpMiddleware) HandlerFunc(h http.HandlerFunc) http.HandlerFunc {
	if wr.tracer == nil {
		return h
	}

	return func(w http.ResponseWriter, r *http.Request) {
		operationName := wr.opNameFunc(r)

		ctx, _ := wr.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		sp := wr.tracer.StartSpan(operationName, ext.RPCServerOption(ctx))
		ext.HTTPMethod.Set(sp, r.Method)

		sct := &statusCodeTracker{ResponseWriter: w}
		r = r.WithContext(opentracing.ContextWithSpan(r.Context(), sp))

		defer func() {
			ext.HTTPStatusCode.Set(sp, uint16(sct.status))
			if sct.status >= http.StatusInternalServerError || !sct.wroteheader {
				ext.Error.Set(sp, true)
			}
			sp.Finish()
		}()

		h(sct.wrappedResponseWriter(), r)
	}
}
