package pprofwrapper

import (
	"net/http"

	"github.com/gorilla/mux"
)

type routerHelper struct {
	router *mux.Router
}

func (r *routerHelper) Handle(path string, handler http.Handler) {
	_ = r.router.Handle(path, handler)
}

func (r *routerHelper) HandleFunc(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	_ = r.router.HandleFunc(path, handler)
}
