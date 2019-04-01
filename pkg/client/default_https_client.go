package client

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
)

func DefaultHTTPSClient(certPool *x509.CertPool) HTTPSClient {
	transport := http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{
		RootCAs: certPool,
	}

	return &http.Client{Transport: transport}
}
