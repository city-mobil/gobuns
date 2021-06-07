package external

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type exSuite struct {
	suite.Suite

	correctSrv *httptest.Server
	timeoutSrv *httptest.Server
}

func newExSuite() *exSuite {
	correctSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"Hello, world!"}`))
	}))

	timeoutSrv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		time.Sleep(100 * time.Millisecond)
	}))

	return &exSuite{
		correctSrv: correctSrv,
		timeoutSrv: timeoutSrv,
	}
}

func TestNewClientWithTLS(t *testing.T) {
	cfg := Config{}
	cfg = cfg.withDefaults()
	cfg.PublicCertPath = "testdata/certs/client.crt"
	cfg.PrivateCertPath = "testdata/certs/client.key"

	client, err := New(&cfg)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestNewClientWithInvalidTLS(t *testing.T) {
	cfg := Config{}
	cfg = cfg.withDefaults()
	cfg.PublicCertPath = "unknown/client.crt"
	cfg.PrivateCertPath = "unknown/client.key"

	client, err := New(&cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClient_EmptyName(t *testing.T) {
	cfg := Config{}
	cfg = cfg.withDefaults()
	cfg.Name = ""
	cfg.Metrics.Collect = true

	client, err := New(&cfg)
	assert.Equal(t, ErrEmptyName, err)
	assert.Nil(t, client)
}

func TestExternal(t *testing.T) {
	suite.Run(t, newExSuite())
}

func (s *exSuite) TearDownSuite() {
	s.correctSrv.Close()
	s.timeoutSrv.Close()
}

func (s *exSuite) TestGet() {
	t := s.T()

	var tests = []struct {
		name    string
		method  string
		srvURL  string
		wantErr bool
	}{
		{
			name:    "InvalidMethod",
			method:  http.MethodPost,
			srvURL:  "https://google.com",
			wantErr: true,
		},
		{
			name:    "RequestTimeout",
			method:  http.MethodGet,
			srvURL:  s.timeoutSrv.URL,
			wantErr: true,
		},
		{
			name:   "SuccessfulRequest",
			method: http.MethodGet,
			srvURL: s.correctSrv.URL,
		},
	}

	cfg := NewConfig("pref")()
	cfg.RequestTimeout = 50 * time.Millisecond
	cl, err := New(cfg)
	require.NoError(t, err)

	for _, tt := range tests {
		v := tt
		t.Run(v.name, func(t *testing.T) {
			req, err := http.NewRequest(v.method, v.srvURL, nil)
			require.NoError(t, err)

			resp, err := cl.Get(context.Background(), req)
			defer func() {
				if resp != nil {
					_ = resp.Body.Close()
				}
			}()
			if v.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				if assert.NotNil(t, resp) {
					assert.Equal(t, http.StatusOK, resp.StatusCode)
				}
			}
		})
	}
}

func (s *exSuite) TestBrokenURL() {
	t := s.T()

	cfg := Config{}
	cfg = cfg.withDefaults()
	cl, err := New(&cfg)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, s.correctSrv.URL, nil)
	require.NoError(t, err)
	req.URL = nil

	_, err = cl.Get(context.Background(), req) // nolint:bodyclose
	assert.Error(t, err)
}

func (s *exSuite) TestPost() {
	t := s.T()

	var tests = []struct {
		name    string
		method  string
		srvURL  string
		wantErr bool
	}{
		{
			name:    "InvalidMethod",
			method:  http.MethodGet,
			srvURL:  "https://google.com",
			wantErr: true,
		},
		{
			name:    "RequestTimeout",
			method:  http.MethodPost,
			srvURL:  s.timeoutSrv.URL,
			wantErr: true,
		},
		{
			name:   "SuccessfulRequest",
			method: http.MethodPost,
			srvURL: s.correctSrv.URL,
		},
	}

	cfg := Config{}
	cfg = cfg.withDefaults()
	cfg.RequestTimeout = 50 * time.Millisecond
	cl, err := New(&cfg)
	require.NoError(t, err)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.srvURL, nil)
			require.NoError(t, err)

			resp, err := cl.Post(context.Background(), req)
			defer func() {
				if resp != nil {
					_ = resp.Body.Close()
				}
			}()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				if assert.NotNil(t, resp) {
					assert.Equal(t, http.StatusOK, resp.StatusCode)
				}
			}
		})
	}
}

func (s *exSuite) TestDo() {
	t := s.T()

	span := opentracing.StartSpan("test")
	spanCtx := opentracing.ContextWithSpan(context.Background(), span)

	var tests = []struct {
		name    string
		ctx     context.Context
		srvURL  string
		wantErr bool
	}{
		{
			name:    "RequestTimeout",
			ctx:     context.Background(),
			srvURL:  s.timeoutSrv.URL,
			wantErr: true,
		},
		{
			name:   "SuccessfulRequest",
			ctx:    context.Background(),
			srvURL: s.correctSrv.URL,
		},
		{
			name:    "RequestTimeoutWithOT",
			ctx:     spanCtx,
			srvURL:  s.timeoutSrv.URL,
			wantErr: true,
		},
		{
			name:   "SuccessfulRequestWithOT",
			ctx:    spanCtx,
			srvURL: s.correctSrv.URL,
		},
	}

	cfg := Config{}
	cfg = cfg.withDefaults()
	cfg.RequestTimeout = 50 * time.Millisecond
	cl, err := New(&cfg)
	require.NoError(t, err)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tt.srvURL, nil)
			require.NoError(t, err)

			resp, err := cl.Do(tt.ctx, req)
			defer func() {
				if resp != nil {
					_ = resp.Body.Close()
				}
			}()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				if assert.NotNil(t, resp) {
					assert.Equal(t, http.StatusOK, resp.StatusCode)
				}
			}
		})
	}
}

func setupBrokenReaderServer(t *testing.T, maxBytesRead int, maxTimeout time.Duration) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, maxBytesRead)
		_, err := r.Body.Read(body)
		assert.NoError(t, err)
		time.Sleep(maxTimeout)
	}))
	return srv
}

func (s *exSuite) TestPostWithRetries() {
	t := s.T()
	brokenReaderSrv := setupBrokenReaderServer(t, 8, 100*time.Millisecond)
	defer brokenReaderSrv.Close()

	var tests = []struct {
		name    string
		srvURL  string
		body    []byte
		wantErr bool
	}{
		{
			name:   "SuccessfulRequest",
			srvURL: s.correctSrv.URL,
			body:   []byte("test"),
		},
		{
			name:    "InvalidDataLength",
			srvURL:  brokenReaderSrv.URL,
			body:    []byte("some long body..."),
			wantErr: true,
		},
	}

	cfg := NewConfig("post_with_retries")()
	cfg.RequestTimeout = 50 * time.Millisecond
	cl, err := New(cfg)
	require.NoError(t, err)

	for _, v := range tests {
		v := v
		t.Run(v.name, func(t *testing.T) {
			rdr := bytes.NewReader(v.body)
			req, err := http.NewRequest(http.MethodPost, v.srvURL, rdr)
			require.NoError(t, err)

			resp, err := cl.Post(context.Background(), req)
			defer func() {
				if resp != nil {
					_ = resp.Body.Close()
				}
			}()

			if v.wantErr {
				if assert.Error(t, err) {
					assert.NotRegexp(t, regexp.MustCompile("with Body length"), err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func (s *exSuite) TestDefaultClient() {
	t := s.T()

	brokenReaderSrv := setupBrokenReaderServer(t, 8, defRequestTimeout+50*time.Millisecond)
	defer brokenReaderSrv.Close()

	var tests = []struct {
		name    string
		srvURL  string
		body    []byte
		wantErr bool
	}{
		{
			name:   "SuccessfulRequest",
			srvURL: s.correctSrv.URL,
			body:   []byte("test"),
		},
		{
			name:    "InvalidDataLength",
			srvURL:  brokenReaderSrv.URL,
			body:    []byte("some long body..."),
			wantErr: true,
		},
	}

	cl := DefaultClient

	for _, v := range tests {
		v := v
		t.Run(v.name, func(t *testing.T) {
			rdr := bytes.NewReader(v.body)
			req, err := http.NewRequest(http.MethodPost, v.srvURL, rdr)
			require.NoError(t, err)

			resp, err := cl.Post(context.Background(), req)
			defer func() {
				if resp != nil {
					_ = resp.Body.Close()
				}
			}()

			if v.wantErr {
				if assert.Error(t, err) {
					assert.NotRegexp(t, regexp.MustCompile("with Body length"), err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
