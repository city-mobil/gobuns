package httputil

import (
	"net"
	"net/http"
	"strings"
)

type IPLookup struct {
	// Places is a list of places to look up IP address.
	// Default is "RemoteAddr", "X-Real-IP", "X-Forwarded-For".
	// You can rearrange the order as you like.
	Places []string

	// ForwardedForIndex is an index of item
	// from header X-Forwarded-For which should be treated as client IP.
	ForwardedForIndex int
}

func NewIPLookup() *IPLookup {
	return &IPLookup{
		Places:            []string{"RemoteAddr", "X-Real-IP", "X-Forwarded-For"},
		ForwardedForIndex: 0,
	}
}

func NewCustomIPLookup(places []string, forwardedForIndex int) *IPLookup {
	return &IPLookup{
		Places:            places,
		ForwardedForIndex: forwardedForIndex,
	}
}

// GetRemoteIP returns IP from HTTP request headers or an empty string if nothing found.
func (ipl *IPLookup) GetRemoteIP(r *http.Request) string {
	for _, lookup := range ipl.Places {
		if lookup == "RemoteAddr" {
			// 1. Cover the basic use cases for both ipv4 and ipv6
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				// 2. Upon error, just return the remote addr.
				return r.RemoteAddr
			}
			return ip
		}

		forwardedFor := r.Header.Get("X-Forwarded-For")
		if lookup == "X-Forwarded-For" && forwardedFor != "" {
			// X-Forwarded-For is potentially a list of addresses separated with ","
			parts := strings.Split(forwardedFor, ",")

			partIndex := ipl.ForwardedForIndex
			if partIndex >= len(parts) {
				partIndex = len(parts) - 1
			}

			return strings.TrimSpace(parts[partIndex])
		}

		realIP := r.Header.Get("X-Real-IP")
		if lookup == "X-Real-IP" && realIP != "" {
			return realIP
		}
	}

	return ""
}
