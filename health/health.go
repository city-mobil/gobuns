// Package health contains API for working with health checks described
// in https://inadarei.github.io/rfc-healthcheck
package health

import (
	"encoding/json"
	"net/http"

	"github.com/city-mobil/gobuns/promlib"
)

type CheckStatus string

const (
	CheckStatusPass CheckStatus = "pass"
	CheckStatusWarn CheckStatus = "warn"
	CheckStatusFail CheckStatus = "fail"
)

// WarnError is a error which sets 'warn' status for some not-passed health/alive check.
type WarnError struct {
	Message string
}

func (e *WarnError) Error() string {
	return e.Message
}

// FailError is a error which sets 'warn' status for some not-passed health/alive check.
//
// All the unknown errors which can not be casted to WarnError or FailError are threated as
// errors, leading to fail of the whole health/alive check.
type FailError struct {
	Message string
}

func (e *FailError) Error() string {
	return e.Message
}

type Checks map[string]*CheckResult

type CheckResult struct {
	ComponentID       string      `json:"component_id"`
	ComponentType     string      `json:"component_type"`
	ObservedValue     interface{} `json:"observed_value"`
	ObservedUnit      string      `json:"observed_unit,omitempty"`
	Status            CheckStatus `json:"status"`
	Output            string      `json:"output"`
	AffectedEndpoints []string    `json:"affected_endpoints,omitempty"`
	Error             error       `json:"-"`
}

type CheckResponse struct {
	Status      CheckStatus `json:"status"`
	Version     string      `json:"version,omitempty"`
	ReleaseID   string      `json:"release_id,omitempty"`
	ServiceID   string      `json:"service_id,omitempty"`
	Description string      `json:"description,omitempty"`
	Output      string      `json:"output,omitempty"`
	Checks      Checks      `json:"checks,omitempty"`
	// TODO(a.petrukhin): add about
}

func NewHandler(ch Checker, handlerName string) func(w http.ResponseWriter, r *http.Request) {
	handler := promlib.NewMiddleware(promlib.DefHTTPRequestDurBuckets, promlib.WithHistogramName(handlerName))
	return handler.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := ch.CheckContext(r.Context())
		data, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if res.Status != CheckStatusFail {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, _ = w.Write(data)
	})
}
