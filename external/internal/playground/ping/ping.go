package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	jaegercfg "github.com/uber/jaeger-client-go/config"

	"github.com/city-mobil/gobuns/config"
	"github.com/city-mobil/gobuns/external"
	"github.com/city-mobil/gobuns/promlib"
)

func main() {
	cfgFn := external.NewConfig("ping")

	err := config.InitOnce()
	if err != nil {
		log.Fatalf("failed to init config: %s", err)
	}

	jaegerCfg, err := jaegercfg.FromEnv()
	if err != nil {
		log.Fatalf("failed to read jaeger configuration: %s", err)
	}
	closer, err := jaegerCfg.InitGlobalTracer("ping")
	if err != nil {
		log.Fatalf("failed to init Open Tracing: %s", err)
	}
	defer closer.Close()

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:    ":9091",
		Handler: router,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalf("listen error: %s", err)
		}
	}()

	externalCfg := cfgFn()
	externalCfg.Metrics.Options = []promlib.InstrumentOption{
		promlib.InstrumentWithPath(func(r *http.Request) string {
			return r.Host
		}),
	}
	client, err := external.New(externalCfg)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodGet, "http://pong:9090/ping", nil)
	if err != nil {
		log.Fatalf("failed to create a request: %s", err)
	}

	go func() {
		for {
			span, ctx := opentracing.StartSpanFromContext(context.Background(), "GetPongOp")

			q := req.URL.Query()
			q.Set("delay", fmt.Sprintf("%dms", rand.Intn(100)))
			req.URL.RawQuery = q.Encode()

			resp, err := client.Get(ctx, req)
			log.Printf("request error: %v", err)
			if resp != nil {
				log.Printf("response status: %v", resp.StatusCode)
				data, _ := ioutil.ReadAll(resp.Body)
				log.Printf("response: %v", string(data))

				_ = resp.Body.Close()
			}

			span.Finish()
			time.Sleep(500 * time.Millisecond)
		}
	}()

	notify := make(chan os.Signal, 1)
	signal.Notify(notify, syscall.SIGINT, syscall.SIGTERM)

	<-notify

	err = srv.Close()
	if err != nil {
		log.Fatalf("failed to gracefully shutdown: %s", err)
	}
}
