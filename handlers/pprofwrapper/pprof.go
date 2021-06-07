// Package pprofwrapper contains wrappers for golang pprof http profiler.
package pprofwrapper

import (
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/mux"
)

type Mux interface {
	Handle(string, http.Handler)
	HandleFunc(string, func(http.ResponseWriter, *http.Request))
}

func register(cfg *Config, mx Mux) {
	mx.HandleFunc("/debug/pprof", pprof.Index)
	mx.Handle("/debug/allocs", pprof.Handler("allocs"))
	mx.Handle("/debug/block", pprof.Handler("block"))

	if cfg.CmdlineEnabled {
		mx.Handle("/debug/cmdline", pprof.Handler("cmdline"))
	}

	mx.Handle("/debug/goroutine", pprof.Handler("goroutine"))
	mx.Handle("/debug/heap", pprof.Handler("heap"))
	mx.Handle("/debug/mutex", pprof.Handler("mutex"))
	mx.HandleFunc("/debug/profile", pprof.Profile)
	mx.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mx.Handle("/debug/threadcreate", pprof.Handler("threadcreate"))
	mx.Handle("/debug/trace", pprof.Handler("trace"))
}

// RegisterDefaultMux registers default http.ServeMux for pprof usage.
//
// Note, that DefaultServeMux MUST BE RESET before usage.
func RegisterDefaultMux(cfg *Config) {
	register(cfg, http.DefaultServeMux)
}

// RegisterRouter registers given router for pprof usage.
func RegisterRouter(cfg *Config, router *mux.Router) {
	helper := &routerHelper{
		router: router,
	}
	register(cfg, helper)
}

// NewHandler creates new handler for pprof usage.
func NewHandler(cfg *Config) http.Handler {
	mx := http.NewServeMux()
	register(cfg, mx)

	return mx
}

// NewServer creates new http server for pprof usage only.
func NewServer(cfg *Config) *http.Server {
	handler := NewHandler(cfg)

	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: handler,
	}

	return srv
}
