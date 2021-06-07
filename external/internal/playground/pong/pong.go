package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	pingHd := func(w http.ResponseWriter, req *http.Request) {
		delay := req.URL.Query().Get("delay")
		if delay != "" {
			dur, err := time.ParseDuration(delay)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = io.WriteString(w, fmt.Sprintf("wrong delay duration: %s", err))

				return
			}

			time.Sleep(dur)
		}

		_, _ = io.WriteString(w, "pong")
	}
	http.HandleFunc("/ping", pingHd)

	srv := &http.Server{
		Addr: ":9090",
	}

	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalf("listen error: %s", err)
		}
	}()

	notify := make(chan os.Signal, 1)
	signal.Notify(notify, syscall.SIGINT, syscall.SIGTERM)

	<-notify

	err := srv.Close()
	if err != nil {
		log.Fatalf("failed to gracefully shutdown: %s", err)
	}
}
