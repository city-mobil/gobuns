package promlib

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func TestInstrumentRoundTripper(t *testing.T) {
	client := http.DefaultClient
	client.Timeout = 1 * time.Second

	reg := prometheus.NewRegistry()
	opts := []InstrumentOption{
		InstrumentWithPath(func(r *http.Request) string {
			return r.URL.Host + r.URL.Path
		}),
		InstrumentWithRegisterer(reg),
	}
	client.Transport = InstrumentRoundTripper("test", http.DefaultTransport, opts...)
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	resp, err := client.Get(backend.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	mfs, err := reg.Gather()
	if err != nil {
		t.Fatal(err)
	}
	if want, got := 3, len(mfs); want != got {
		t.Fatalf("unexpected number of metric families gathered, want %d, got %d", want, got)
	}
	for _, mf := range mfs {
		if len(mf.Metric) == 0 {
			t.Errorf("metric family %s must not be empty", mf.GetName())
		}
	}
}
