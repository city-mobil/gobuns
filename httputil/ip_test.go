package httputil

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewIPLookup(t *testing.T) {
	lookup := NewIPLookup()
	assert.Empty(t, lookup.ForwardedForIndex)
	assert.Equal(t, []string{"RemoteAddr", "X-Real-IP", "X-Forwarded-For"}, lookup.Places)
}

func TestNewCustomIPLookup(t *testing.T) {
	ipLookups := []string{"X-Real-IP", "RemoteAddr", "X-IP-Header"}
	forwardedForIndex := 2
	lookup := NewCustomIPLookup(ipLookups, forwardedForIndex)
	assert.Equal(t, forwardedForIndex, lookup.ForwardedForIndex)
	assert.Equal(t, ipLookups, lookup.Places)
}

func TestIPLookup_GetRemoteIP_EmptyLookup(t *testing.T) {
	lookup := NewCustomIPLookup(nil, 0)

	request, err := http.NewRequest("GET", "/", strings.NewReader("Hello, world!"))
	require.Nil(t, err)

	request.Header.Set("X-Real-IP", "10.0.0.10")

	ip := lookup.GetRemoteIP(request)
	assert.Empty(t, ip)
}

func TestIPLookup_GetRemoteIP(t *testing.T) {
	ipLookups := []string{"RemoteAddr", "X-Real-IP"}
	lookup := NewCustomIPLookup(ipLookups, 0)
	ipv6 := "2601:7:1c82:4097:59a0:a80b:2841:b8c8"

	request, err := http.NewRequest("GET", "/", strings.NewReader("Hello, world!"))
	require.Nil(t, err)

	request.Header.Set("X-Real-IP", ipv6)

	ip := lookup.GetRemoteIP(request)
	assert.Equal(t, request.RemoteAddr, ip)
	assert.NotEqual(t, ipv6, ip, "X-Real-IP should have been skipped")
}

func TestIPLookup_GetRemoteIP_RemoteAddr(t *testing.T) {
	ipLookups := []string{"RemoteAddr"}
	lookup := NewCustomIPLookup(ipLookups, 0)

	request := httptest.NewRequest("GET", "/", strings.NewReader("Hello, world!"))

	ip := lookup.GetRemoteIP(request)
	assert.Equal(t, "192.0.2.1", ip)
}

func TestIPLookup_GetRemoteIP_ForwardedFor(t *testing.T) {
	ipLookups := []string{"X-Forwarded-For", "X-Real-IP", "RemoteAddr"}
	lookup := NewCustomIPLookup(ipLookups, 0)
	ipv4 := "10.10.10.11"
	ipv6 := "2601:7:1c82:4097:59a0:a80b:2841:b8c9"

	request, err := http.NewRequest("GET", "/", strings.NewReader("Hello, world!"))
	require.Nil(t, err)

	request.Header.Set("X-Forwarded-For", ipv4)
	request.Header.Set("X-Real-IP", ipv6)

	ip := lookup.GetRemoteIP(request)
	assert.Equal(t, ipv4, ip)
	assert.NotEqual(t, ipv6, ip, "X-Real-IP should have been skipped")
}

func TestIPLookup_GetRemoteIP_RealIP(t *testing.T) {
	ipLookups := []string{"X-Real-IP", "X-Forwarded-For", "RemoteAddr"}
	lookup := NewCustomIPLookup(ipLookups, 0)
	ipv4 := "10.10.10.12"
	ipv6 := "2601:7:1c82:4097:59a0:a80b:2841:b8c7"

	request, err := http.NewRequest("GET", "/", strings.NewReader("Hello, world!"))
	require.Nil(t, err)

	request.Header.Set("X-Forwarded-For", ipv4)
	request.Header.Set("X-Real-IP", ipv6)

	ip := lookup.GetRemoteIP(request)
	assert.Equal(t, ipv6, ip)
	assert.NotEqual(t, ipv4, ip, "X-Forwarded-For should have been skipped")
}

func TestIPLookup_GetRemoteIP_MultipleForwardedFor(t *testing.T) {
	ipLookups := []string{"X-Forwarded-For", "X-Real-IP", "RemoteAddr"}
	lookup := NewCustomIPLookup(ipLookups, 0)
	ipv6 := "2601:7:1c82:4097:59a0:a80b:2841:b8c8"

	request, err := http.NewRequest("GET", "/", strings.NewReader("Hello, world!"))
	require.Nil(t, err)

	request.Header.Set("X-Real-IP", ipv6)

	// Missing X-Forwarded-For should not break things.
	ip := lookup.GetRemoteIP(request)
	assert.Equal(t, ipv6, ip, "X-Real-IP should have been chosen because X-Forwarded-For is missing")

	request.Header.Set("X-Forwarded-For", "10.10.10.10,10.10.10.11")

	// Should get the first one
	ip = lookup.GetRemoteIP(request)
	assert.Equal(t, "10.10.10.10", ip)
	assert.NotEqual(t, ipv6, ip, "X-Real-IP should have been skipped")

	// Should get the last
	lookup.ForwardedForIndex = 1
	ip = lookup.GetRemoteIP(request)
	assert.Equal(t, "10.10.10.11", ip)
	assert.NotEqual(t, ipv6, ip, "X-Real-IP should have been skipped")

	// What about index out of bound? GetRemoteIP should simply choose the last one.
	lookup.ForwardedForIndex = 2
	ip = lookup.GetRemoteIP(request)
	assert.Equal(t, "10.10.10.11", ip)
	assert.NotEqual(t, ipv6, ip, "X-Real-IP should have been skipped")
}
