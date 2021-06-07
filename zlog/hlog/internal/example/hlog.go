package main

import (
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/city-mobil/gobuns/zlog"
	"github.com/city-mobil/gobuns/zlog/hlog"
)

func final(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)
	log.Info().Str("status", "ok").Msg("request finished")

	_, _ = w.Write([]byte("OK"))
}

func main() {
	log := zlog.New(os.Stdout)
	mux := http.NewServeMux()

	fh := http.HandlerFunc(final)
	mw := hlog.RequestIDHandler("request_id", hlog.JaegerTraceHeaderName, true)(
		hlog.RemoteAddrHandler("remote")(
			hlog.MethodHandler("method")(fh),
		),
	)

	mux.Handle("/", hlog.NewHandler(log)(mw))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1:80"
	req.Header = http.Header{
		hlog.JaegerTraceHeaderName: []string{"514bbe5bb5251c92bd07a9846f4a1ab6"},
	}
	mux.ServeHTTP(httptest.NewRecorder(), req)
}

// Output: {"level":"info","request_id":"514bbe5bb5251c92bd07a9846f4a1ab6","remote":"127.0.0.1","method":"GET","status":"ok","time":"2021-05-31T15:46:19+03:00","message":"request finished"}
