package promlib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLNumberNormalizer(t *testing.T) {
	tests := []struct {
		name      string
		req       *http.Request
		templates []string
		want      string
	}{
		{
			name: "NoTemplates",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/api/v0/order/100500/active", nil)
			}(),
			templates: nil,
			want:      "/api/v0/order/100500/active",
		},
		{
			name: "NumberPartNotFound",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/api/v0/order/active", nil)
			}(),
			templates: []string{":id"},
			want:      "/api/v0/order/active",
		},
		{
			name: "NumberPartInTheMiddle",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/api/v0/order/100500/active", nil)
			}(),
			templates: []string{":id"},
			want:      "/api/v0/order/:id/active",
		},
		{
			name: "NumberPartAtTheEnd",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/api/v0/order/100500", nil)
			}(),
			templates: []string{":id"},
			want:      "/api/v0/order/:id",
		},
		{
			name: "MultipleNumberParts",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/api/v0/department/100500/user/12", nil)
			}(),
			templates: []string{":did", ":uid"},
			want:      "/api/v0/department/:did/user/:uid",
		},
		{
			name: "NotEnoughTemplates",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/api/v0/department/100500/user/12", nil)
			}(),
			templates: []string{":id"},
			want:      "/api/v0/department/:id/user/:id",
		},
		{
			name: "RequestWithQuery",
			req: func() *http.Request {
				req := httptest.NewRequest("GET", "/api/v0/department/100500", nil)
				q := req.URL.Query()
				q.Set("a", "12")
				req.URL.RawQuery = q.Encode()
				return req
			}(),
			templates: []string{":id"},
			want:      "/api/v0/department/:id",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := URLNumberNormalizer(tt.req, tt.templates)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestURLHEXNumberNormalizer(t *testing.T) {
	tests := []struct {
		name      string
		req       *http.Request
		templates []string
		want      string
	}{
		{
			name: "NoTemplates",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/api/v0/order/8c86/active", nil)
			}(),
			templates: nil,
			want:      "/api/v0/order/8c86/active",
		},
		{
			name: "HexPartNotFound",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/api/v0/order/active", nil)
			}(),
			templates: []string{":id"},
			want:      "/api/v0/order/active",
		},
		{
			name: "HexPartInTheMiddle",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/api/v0/order/18c86399caf0da8aed/active", nil)
			}(),
			templates: []string{":id"},
			want:      "/api/v0/order/:id/active",
		},
		{
			name: "HexPartAtTheEnd",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/api/v0/order/18c86399", nil)
			}(),
			templates: []string{":id"},
			want:      "/api/v0/order/:id",
		},
		{
			name: "MultipleHexParts",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/api/v0/department/18c86399/user/12", nil)
			}(),
			templates: []string{":did", ":uid"},
			want:      "/api/v0/department/:did/user/:uid",
		},
		{
			name: "NotEnoughTemplates",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/api/v0/department/100500/user/12", nil)
			}(),
			templates: []string{":id"},
			want:      "/api/v0/department/:id/user/:id",
		},
		{
			name: "RequestWithQuery",
			req: func() *http.Request {
				req := httptest.NewRequest("GET", "/api/v0/department/100500", nil)
				q := req.URL.Query()
				q.Set("a", "12")
				req.URL.RawQuery = q.Encode()
				return req
			}(),
			templates: []string{":id"},
			want:      "/api/v0/department/:id",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := URLHEXNumberNormalizer(tt.req, tt.templates)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestURLIndicesNormalizer(t *testing.T) {
	tests := []struct {
		name      string
		req       *http.Request
		templates []string
		indices   []int
		want      string
	}{
		{
			name: "NoTemplates",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/api/v0/order/8c86/active", nil)
			}(),
			templates: nil,
			indices:   nil,
			want:      "/api/v0/order/8c86/active",
		},
		{
			name: "NoIndices",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/api/v0/order/active", nil)
			}(),
			templates: []string{":id"},
			indices:   nil,
			want:      "/api/v0/order/active",
		},
		{
			name: "IndexInTheMiddle",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/api/v0/order/18c86399caf0da8aed/active", nil)
			}(),
			templates: []string{":id"},
			indices:   []int{3},
			want:      "/api/v0/order/:id/active",
		},
		{
			name: "IndexAtTheEnd",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/api/v0/order/18c86399", nil)
			}(),
			templates: []string{":id"},
			indices:   []int{3},
			want:      "/api/v0/order/:id",
		},
		{
			name: "MultipleIndexes",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/api/v0/department/18c86399/user/12", nil)
			}(),
			templates: []string{":did", ":uid"},
			indices:   []int{3, 5},
			want:      "/api/v0/department/:did/user/:uid",
		},
		{
			name: "NotEnoughTemplates",
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/api/v0/department/100500/user/12", nil)
			}(),
			templates: []string{":id"},
			indices:   []int{3, 5},
			want:      "/api/v0/department/:id/user/:id",
		},
		{
			name: "RequestWithQuery",
			req: func() *http.Request {
				req := httptest.NewRequest("GET", "/api/v0/department/100500", nil)
				q := req.URL.Query()
				q.Set("a", "12")
				req.URL.RawQuery = q.Encode()
				return req
			}(),
			templates: []string{":id"},
			indices:   []int{2},
			want:      "/api/v0/:id/100500",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := URLIndicesNormalizer(tt.req, tt.templates, tt.indices)
			assert.Equal(t, tt.want, got)
		})
	}
}

func BenchmarkParse(b *testing.B) {
	req := httptest.NewRequest("GET", "/api/v0/order/100500/active", nil)
	templates := []string{":id"}

	b.ReportAllocs()
	b.ResetTimer()
	var res string
	for i := 0; i < b.N; i++ {
		res = URLNumberNormalizer(req, templates)
	}
	_, _ = fmt.Fprint(ioutil.Discard, res)
}
