package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/city-mobil/gobuns/config"
	"github.com/city-mobil/gobuns/graceful"
	"github.com/city-mobil/gobuns/health"
	"github.com/city-mobil/gobuns/kafka"
	"github.com/city-mobil/gobuns/zlog"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	kf "github.com/segmentio/kafka-go"
)

func main() {
	cfg := kafka.NewProducerConfig("")
	err := config.InitOnce()
	if err != nil {
		log.Fatal(err)
	}

	logger := zlog.New(os.Stdout)

	producer := kafka.NewAsyncProducer(logger, cfg())
	producer.SetCompletionCallback(func(msgs []kf.Message, err error) {
		if err == nil {
			return
		}

		// NOTE(a.petrukhin): например, попытаемся отправить ещё раз.
		_ = producer.Produce(context.Background(), msgs...)
	})

	ch := health.NewChecker(health.CheckerOptions{
		ReleaseID: "1",
		ServiceID: "async_producer",
		Version:   "v1.0.0",
	})
	ch.AddCallback("kafka:producer", kafka.NewProducerHealthCheckCallback(producer))

	http.HandleFunc("/health", health.NewHandler(ch, "health"))
	http.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr: ":4242",
	}
	go func() {
		_ = srv.ListenAndServe()
	}()

	go func() {
		for {
			err = producer.Produce(context.Background(), []kf.Message{
				{
					Topic: "orders",
					Value: []byte("some_value"),
				},
			}...)
			if err != nil {
				logger.Err(err).Send()
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	_ = graceful.WaitShutdown()
}
