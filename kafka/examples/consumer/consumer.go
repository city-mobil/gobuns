package consumer

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/city-mobil/gobuns/graceful"

	"github.com/city-mobil/gobuns/health"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/city-mobil/gobuns/zlog"

	"github.com/city-mobil/gobuns/config"
	"github.com/city-mobil/gobuns/kafka"
)

func main() {
	cfg := kafka.NewConsumerConfig("")
	err := config.InitOnce()
	if err != nil {
		log.Fatal(err)
	}

	logger := zlog.New(os.Stdout)
	consumer := kafka.NewConsumer(logger, cfg())
	ch := health.NewChecker(health.CheckerOptions{
		ReleaseID: "1",
		ServiceID: "consumer",
		Version:   "v1.0.0",
	})
	ch.AddCallback("kafka:consumer", kafka.NewConsumerHealthCheckCallback(consumer))

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
			msg, err := consumer.FetchMessage(context.Background())
			if err != nil {
				logger.Err(err).Send()
			}
			logger.Info().Bytes("message", msg.Value).Msg("got message")
			time.Sleep(10 * time.Millisecond)
		}
	}()

	_ = graceful.WaitShutdown()
}
