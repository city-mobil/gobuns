package promlib

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"

	"github.com/prometheus/client_golang/prometheus"
)

type httpSuite struct {
	suite.Suite
}

func TestHTTPSuite(t *testing.T) {
	suite.Run(t, new(httpSuite))
}

func (s *httpSuite) SetupTest() {
	reg := prometheus.NewRegistry()
	prometheus.DefaultGatherer = reg
	prometheus.DefaultRegisterer = reg
}

func (s *httpSuite) TestHTTPMiddleware_WithCustomPath() {
	t := s.T()

	mw := NewMiddleware([]float64{0.5, 1}, WithCustomPath(func(r *http.Request) string {
		return URLNumberNormalizer(r, []string{":id"})
	}))

	req := httptest.NewRequest("GET", "/api/v0/users/1", nil)

	rr := httptest.NewRecorder()
	mw.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}).ServeHTTP(rr, req)

	metrics := dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, "# HELP http_request_duration_seconds Duration of HTTP request")
	assert.Contains(t, metrics, "# TYPE http_request_duration_seconds histogram")
	assert.Contains(t, metrics, `http_request_duration_seconds_bucket{code="200",method="get",path="/api/v0/users/:id",le="0.5"} 1`)
	assert.Contains(t, metrics, `http_request_duration_seconds_bucket{code="200",method="get",path="/api/v0/users/:id",le="1"} 1`)
	assert.Contains(t, metrics, `http_request_duration_seconds_bucket{code="200",method="get",path="/api/v0/users/:id",le="+Inf"} 1`)
	assert.Contains(t, metrics, `http_request_duration_seconds_count{code="200",method="get",path="/api/v0/users/:id"} 1`)
}

func (s *httpSuite) TestHTTPMiddleware_WithUserAgent() {
	t := s.T()

	mw := NewMiddleware([]float64{0.5, 1}, WithUserAgentLabel())

	req := httptest.NewRequest("GET", "/health", nil)
	req.Header.Add("User-Agent", "go-buns")

	rr := httptest.NewRecorder()
	mw.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}).ServeHTTP(rr, req)

	metrics := dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, "# HELP http_request_duration_seconds Duration of HTTP request")
	assert.Contains(t, metrics, "# TYPE http_request_duration_seconds histogram")
	assert.Contains(t, metrics, `http_request_duration_seconds_bucket{agent="go-buns",code="200",method="get",path="/health",le="0.5"} 1`)
	assert.Contains(t, metrics, `http_request_duration_seconds_bucket{agent="go-buns",code="200",method="get",path="/health",le="1"} 1`)
	assert.Contains(t, metrics, `http_request_duration_seconds_bucket{agent="go-buns",code="200",method="get",path="/health",le="+Inf"} 1`)
	assert.Contains(t, metrics, `http_request_duration_seconds_count{agent="go-buns",code="200",method="get",path="/health"} 1`)

	req = httptest.NewRequest("GET", "/health", nil)
	req.Header.Add("User-Agent", "Mozilla")
	mw.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}).ServeHTTP(rr, req)

	metrics = dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, `http_request_duration_seconds_bucket{agent="go-buns",code="200",method="get",path="/health",le="0.5"} 1`)
	assert.Contains(t, metrics, `http_request_duration_seconds_bucket{agent="go-buns",code="200",method="get",path="/health",le="1"} 1`)
	assert.Contains(t, metrics, `http_request_duration_seconds_bucket{agent="go-buns",code="200",method="get",path="/health",le="+Inf"} 1`)
	assert.Contains(t, metrics, `http_request_duration_seconds_count{agent="go-buns",code="200",method="get",path="/health"} 1`)
	assert.Contains(t, metrics, `http_request_duration_seconds_bucket{agent="Mozilla",code="200",method="get",path="/health",le="0.5"} 1`)
	assert.Contains(t, metrics, `http_request_duration_seconds_bucket{agent="Mozilla",code="200",method="get",path="/health",le="1"} 1`)
	assert.Contains(t, metrics, `http_request_duration_seconds_bucket{agent="Mozilla",code="200",method="get",path="/health",le="+Inf"} 1`)
	assert.Contains(t, metrics, `http_request_duration_seconds_count{agent="Mozilla",code="200",method="get",path="/health"} 1`)
}

func (s *httpSuite) TestHTTPMiddleware_WithoutUserAgent() {
	t := s.T()

	mw := NewMiddleware([]float64{0.5, 1})

	req := httptest.NewRequest("GET", "/health", nil)
	req.Header.Add("User-Agent", "go-buns")

	rr := httptest.NewRecorder()
	mw.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}).ServeHTTP(rr, req)

	metrics := dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, "# HELP http_request_duration_seconds Duration of HTTP request")
	assert.Contains(t, metrics, "# TYPE http_request_duration_seconds histogram")
	assert.Contains(t, metrics, `http_request_duration_seconds_bucket{code="200",method="get",path="/health",le="0.5"} 1`)
	assert.Contains(t, metrics, `http_request_duration_seconds_bucket{code="200",method="get",path="/health",le="1"} 1`)
	assert.Contains(t, metrics, `http_request_duration_seconds_bucket{code="200",method="get",path="/health",le="+Inf"} 1`)
	assert.Contains(t, metrics, `http_request_duration_seconds_count{code="200",method="get",path="/health"} 1`)
}

func (s *httpSuite) TestHTTPMiddleware_WithCustomName() {
	t := s.T()

	mw := NewMiddleware([]float64{0.5, 1}, WithHistogramName("my_name"))

	req := httptest.NewRequest("GET", "/health", nil)

	rr := httptest.NewRecorder()
	mw.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}).ServeHTTP(rr, req)

	metrics := dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, "# HELP my_name Duration of HTTP request")
	assert.Contains(t, metrics, "# TYPE my_name histogram")
	assert.Contains(t, metrics, `my_name_bucket{code="200",method="get",path="/health",le="0.5"} 1`)
	assert.Contains(t, metrics, `my_name_bucket{code="200",method="get",path="/health",le="1"} 1`)
	assert.Contains(t, metrics, `my_name_bucket{code="200",method="get",path="/health",le="+Inf"} 1`)
	assert.Contains(t, metrics, `my_name_count{code="200",method="get",path="/health"} 1`)
}
