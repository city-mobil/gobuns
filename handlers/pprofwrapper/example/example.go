package main

import (
	"net/http"

	"github.com/city-mobil/gobuns/handlers/pprofwrapper"
)

func main() {
	// NOTE(a.petrukhin): reset DefaultServeMux in order to register pprof again.
	http.DefaultServeMux = http.NewServeMux()

	pprofwrapper.RegisterDefaultMux(&pprofwrapper.Config{})
	_ = http.ListenAndServe(":80", nil)
}
