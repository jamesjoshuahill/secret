package client_test

import (
	"crypto/x509"

	"github.com/jamesjoshuahill/ciphers/pkg/client"
)

func Example() {
	caCert := `-----BEGIN CERTIFICATE-----
...
-----END CERTIFICATE-----`

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM([]byte(caCert))
	if !ok {
		panic("failed to append CA certificate")
	}

	httpsClient := client.DefaultHTTPSClient(certPool)

	client.New("https://127.0.0.1:8080", httpsClient)
}
