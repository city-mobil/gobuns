package hlog

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/city-mobil/gobuns/zlog"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	log := zlog.New(nil).With().
		Str("foo", "bar").
		Logger()
	lh := NewHandler(log)
	h := lh(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		assert.Equal(t, log, l)
	}))
	h.ServeHTTP(nil, &http.Request{})
}

func TestNewDefHandler(t *testing.T) {
	out := &bytes.Buffer{}
	r := &http.Request{
		Header: http.Header{
			NginxTraceHeaderName: []string{"514bbe5bb5251c92bd07a9846f4a1ab6"},
		},
	}
	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.Log().Msg("")
	})
	h := NewDefHandler(zlog.Raw(out), final)
	h.ServeHTTP(nil, r)
	assert.Equal(t, `{"request_id":"514bbe5bb5251c92bd07a9846f4a1ab6"}`+"\n", out.String())
}

func TestURLHandler(t *testing.T) {
	out := &bytes.Buffer{}
	r := &http.Request{
		URL: &url.URL{Path: "/path", RawQuery: "foo=bar"},
	}
	h := URLHandler("url")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.Log().Msg("")
	}))
	h = NewHandler(zlog.Raw(out))(h)
	h.ServeHTTP(nil, r)
	assert.Equal(t, `{"url":"/path?foo=bar"}`+"\n", out.String())
}

func TestMethodHandler(t *testing.T) {
	out := &bytes.Buffer{}
	r := &http.Request{
		Method: "POST",
	}
	h := MethodHandler("method")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.Log().Msg("")
	}))
	h = NewHandler(zlog.Raw(out))(h)
	h.ServeHTTP(nil, r)
	assert.Equal(t, `{"method":"POST"}`+"\n", out.String())
}

func TestRequestHandler(t *testing.T) {
	out := &bytes.Buffer{}
	r := &http.Request{
		Method: "POST",
		URL:    &url.URL{Path: "/path", RawQuery: "foo=bar"},
	}
	h := RequestHandler("request")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.Log().Msg("")
	}))
	h = NewHandler(zlog.Raw(out))(h)
	h.ServeHTTP(nil, r)
	assert.Equal(t, `{"request":"POST /path?foo=bar"}`+"\n", out.String())
}

func TestRemoteAddrHandler(t *testing.T) {
	out := &bytes.Buffer{}
	r := &http.Request{
		RemoteAddr: "1.2.3.4:1234",
	}
	h := RemoteAddrHandler("ip")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.Log().Msg("")
	}))
	h = NewHandler(zlog.Raw(out))(h)
	h.ServeHTTP(nil, r)
	assert.Equal(t, `{"ip":"1.2.3.4"}`+"\n", out.String())
}

func TestRemoteAddrHandlerIPv6(t *testing.T) {
	out := &bytes.Buffer{}
	r := &http.Request{
		RemoteAddr: "[2001:db8:a0b:12f0::1]:1234",
	}
	h := RemoteAddrHandler("ip")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.Log().Msg("")
	}))
	h = NewHandler(zlog.Raw(out))(h)
	h.ServeHTTP(nil, r)
	assert.Equal(t, `{"ip":"2001:db8:a0b:12f0::1"}`+"\n", out.String())
}

func TestUserAgentHandler(t *testing.T) {
	out := &bytes.Buffer{}
	r := &http.Request{
		Header: http.Header{
			"User-Agent": []string{"some user agent string"},
		},
	}
	h := UserAgentHandler("ua")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.Log().Msg("")
	}))
	h = NewHandler(zlog.Raw(out))(h)
	h.ServeHTTP(nil, r)
	assert.Equal(t, `{"ua":"some user agent string"}`+"\n", out.String())
}

func TestRefererHandler(t *testing.T) {
	out := &bytes.Buffer{}
	r := &http.Request{
		Header: http.Header{
			"Referer": []string{"http://foo.com/bar"},
		},
	}
	h := RefererHandler("referer")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.Log().Msg("")
	}))
	h = NewHandler(zlog.Raw(out))(h)
	h.ServeHTTP(nil, r)
	assert.Equal(t, `{"referer":"http://foo.com/bar"}`+"\n", out.String())
}

func TestRequestIDHandler(t *testing.T) {
	out := &bytes.Buffer{}
	r := &http.Request{
		Header: http.Header{
			RFCTraceHeaderName: []string{"093ea86ad8d1aa65"},
		},
	}
	h := RequestIDHandler("id", RFCTraceHeaderName, true)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.Log().Msg("")
	}))
	h = NewHandler(zlog.Raw(out))(h)
	h.ServeHTTP(httptest.NewRecorder(), r)
	assert.Equal(t, `{"id":"093ea86ad8d1aa65"}`+"\n", out.String())
}

func TestCustomHeaderHandler(t *testing.T) {
	out := &bytes.Buffer{}
	r := &http.Request{
		Header: http.Header{
			"X-Request-Id": []string{"514bbe5bb5251c92bd07a9846f4a1ab6"},
		},
	}
	h := CustomHeaderHandler("reqID", "X-Request-Id")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.Log().Msg("")
	}))
	h = NewHandler(zlog.Raw(out))(h)
	h.ServeHTTP(nil, r)
	assert.Equal(t, `{"reqID":"514bbe5bb5251c92bd07a9846f4a1ab6"}`+"\n", out.String())
}

func TestCombinedHandlers(t *testing.T) {
	out := &bytes.Buffer{}
	r := &http.Request{
		Method: "POST",
		URL:    &url.URL{Path: "/path", RawQuery: "foo=bar"},
	}
	h := MethodHandler("method")(RequestHandler("request")(URLHandler("url")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.Log().Msg("")
	}))))
	h = NewHandler(zlog.Raw(out))(h)
	h.ServeHTTP(nil, r)
	assert.Equal(t, `{"method":"POST","request":"POST /path?foo=bar","url":"/path?foo=bar"}`+"\n", out.String())
}

func BenchmarkHandlers(b *testing.B) {
	r := &http.Request{
		Method: "POST",
		URL:    &url.URL{Path: "/path", RawQuery: "foo=bar"},
	}
	h1 := URLHandler("url")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.Log().Msg("")
	}))
	h2 := MethodHandler("method")(RequestHandler("request")(h1))
	handlers := map[string]http.Handler{
		"Single":           NewHandler(zlog.Raw(ioutil.Discard))(h1),
		"Combined":         NewHandler(zlog.Raw(ioutil.Discard))(h2),
		"SingleDisabled":   NewHandler(zlog.Raw(ioutil.Discard).Level(zlog.Disabled))(h1),
		"CombinedDisabled": NewHandler(zlog.Raw(ioutil.Discard).Level(zlog.Disabled))(h2),
	}
	for name := range handlers {
		h := handlers[name]
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				h.ServeHTTP(nil, r)
			}
		})
	}
}

func BenchmarkDataRace(b *testing.B) {
	log := zlog.Raw(nil).With().
		Str("foo", "bar").
		Logger()
	lh := NewHandler(log)
	h := lh(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.UpdateContext(func(c zlog.Context) zlog.Context {
			return c.Str("bar", "baz")
		})
		l.Log().Msg("")
	}))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.ServeHTTP(nil, &http.Request{})
		}
	})
}
