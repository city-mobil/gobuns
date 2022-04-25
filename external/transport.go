package external

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
)

// newDefaultTransport returns a default transport for external HTTP client.
func newDefaultTransport(cfg *Config) (http.RoundTripper, error) {
	insecureTLS := cfg.NoHTTPS
	tlsConfig := &tls.Config{ //nolint:gosec
		MinVersion: castTLSVersion(cfg.MinVersionTLS),
	}
	if cfg.CACertPath != "" {
		caCert, err := ioutil.ReadFile(cfg.CACertPath)
		if err != nil {
			return nil, err
		}
		tlsConfig.RootCAs = x509.NewCertPool()
		tlsConfig.RootCAs.AppendCertsFromPEM(caCert)
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
