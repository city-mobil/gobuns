package external

import (
	"crypto/tls"
	"net/http"
)

// newDefaultTransport returns a default transport for external HTTP client.
func newDefaultTransport(cfg *Config) (http.RoundTripper, error) {
	insecureTLS := cfg.NoHTTPS
	tlsConfig := &tls.Config{ //nolint:gosec
		MinVersion: castTLSVersion(cfg.MinVersionTLS),
	}
	if cfg.PrivateCertPath != "" || cfg.PublicCertPath != "" {
		insecureTLS = false
		cert, err := tls.LoadX509KeyPair(cfg.PublicCertPath, cfg.PrivateCertPath)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	tlsConfig.InsecureSkipVerify = insecureTLS

	if cfg.ForceInsecureSkipVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	tr := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DialContext:         cfg.DialContext,
		DisableKeepAlives:   false,
		MaxIdleConns:        cfg.MaxIdleConns,
		MaxIdleConnsPerHost: cfg.MaxIdleConns,
		IdleConnTimeout:     cfg.IdleConnTimeout,
		TLSClientConfig:     tlsConfig,
	}

	return tr, nil
}
