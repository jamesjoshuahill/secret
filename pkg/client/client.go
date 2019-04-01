// Package client provides a client for the cipher server.
package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const ciphersResourcePath = "/v1/ciphers"

type client struct {
	baseURL     string
	httpsClient HTTPSClient
}

// New returns a client struct with the given server base URL and HTTPS client.
func New(baseURL string, httpsClient HTTPSClient) *client {
	return &client{
		baseURL:     baseURL,
		httpsClient: httpsClient,
	}
}

func (c *client) do(method, url string, body interface{}) (*http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return c.httpsClient.Do(req)
}
