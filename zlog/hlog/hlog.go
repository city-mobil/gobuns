package hlog

import (
	"net"
	"net/http"

	"github.com/city-mobil/gobuns/zlog"
	"github.com/rs/xid"
)

const (
	NginxTraceHeaderName  = "Request-Id"
	JaegerTraceHeaderName = "Uber-Trace-Id"
	RFCTraceHeaderName    = "Trace-Id"
)

// FromRequest returns the logger from the request's context.
func FromRequest(r *http.Request) zlog.Logger {
	return zlog.FromContext(r.Context())
}

// NewDefHandler returns a default handler with injected logger.
func NewDefHandler(log zlog.Logger, next http.Handler) http.Handler {
	h := RequestIDHandler("request_id", NginxTraceHeaderName, true)(next)
	return NewHandler(log)(h)
}

// NewHandler injects log into requests context.
func NewHandler(log zlog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a copy of the logger
			// to prevent data race when using UpdateContext.
			copied := log.With().Logger()
			r = r.WithContext(zlog.NewContext(r.Context(), copied))
			next.ServeHTTP(w, r)
		})
	}
}

// URLHandler adds the requested URL as a field to the context's logger
// using fieldKey as field key.
func URLHandler(fieldKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := zlog.FromContext(r.Context())
			log.UpdateContext(func(c zlog.Context) zlog.Context {
				return c.Str(fieldKey, r.URL.String())
			})
			next.ServeHTTP(w, r)
		})
	}
}

// MethodHandler adds the request method as a field to the context's logger
// using fieldKey as field key.
func MethodHandler(fieldKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := zlog.FromContext(r.Context())
			log.UpdateContext(func(c zlog.Context) zlog.Context {
				return c.Str(fieldKey, r.Method)
			})
			next.ServeHTTP(w, r)
		})
	}
}

// RequestHandler adds the request method and URL as a field to the context's logger
// using fieldKey as field key.
func RequestHandler(fieldKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := zlog.FromContext(r.Context())
			log.UpdateContext(func(c zlog.Context) zlog.Context {
				return c.Str(fieldKey, r.Method+" "+r.URL.String())
			})
			next.ServeHTTP(w, r)
		})
	}
}

// RemoteAddrHandler adds the request's remote address as a field to the context's logger
// using fieldKey as field key.
func RemoteAddrHandler(fieldKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
				log := zlog.FromContext(r.Context())
				log.UpdateContext(func(c zlog.Context) zlog.Context {
					return c.Str(fieldKey, host)
				})
			}
			next.ServeHTTP(w, r)
		})
	}
}

// UserAgentHandler adds the request's user-agent as a field to the context's logger
// using fieldKey as field key.
func UserAgentHandler(fieldKey string) func(next http.Handler) http.Handler {
	return CustomHeaderHandler(fieldKey, "User-Agent")
}

// RefererHandler adds the request's referer as a field to the context's logger
// using fieldKey as field key.
func RefererHandler(fieldKey string) func(next http.Handler) http.Handler {
	return CustomHeaderHandler(fieldKey, "Referer")
}

// RequestIDHandler adds the request's unique ID from the given header
// as a field to the context's logger using fieldKey as field key.
//
// If there is no header in the request and generateReqID is set,
// identifier will be generated automatically.
func RequestIDHandler(fieldKey, header string, generateReqID bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			val := r.Header.Get(header)
			if val == "" && generateReqID {
				val = xid.New().String()
			}
			if val != "" {
				log := zlog.FromContext(r.Context())
				log.UpdateContext(func(c zlog.Context) zlog.Context {
					return c.Str(fieldKey, val)
				})
			}
			next.ServeHTTP(w, r)
		})
	}
}

// CustomHeaderHandler adds given header from request's header as a field to
// the context's logger using fieldKey as field key.
func CustomHeaderHandler(fieldKey, header string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if val := r.Header.Get(header); val != "" {
				log := zlog.FromContext(r.Context())
				log.UpdateContext(func(c zlog.Context) zlog.Context {
					return c.Str(fieldKey, val)
				})
			}
			next.ServeHTTP(w, r)
		})
	}
}
