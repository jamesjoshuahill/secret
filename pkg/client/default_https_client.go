package client

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
)

// DefaultHTTPSClient returns the http.Client struct with the default
// transport and the given certificate pool.
//
// The http.Client struct implements the HTTPSClient interface.
func DefaultHTTPSClient(certPool *x509.CertPool) *http.Client {
	transport := http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{
		RootCAs: certPool,
	}

	return &http.Client{Transport: transport}
}
