// Package client provides a Client for the cipher server.
package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const ciphersResourcePath = "/v1/ciphers"

type Client struct {
	baseURL     string
	httpsClient HTTPSClient
}

// New returns a Client struct with the given baseURL of the server and HTTPSClient.
func New(baseURL string, httpsClient HTTPSClient) *Client {
	return &Client{
		baseURL:     baseURL,
		httpsClient: httpsClient,
	}
}

func (c *Client) do(method, url string, body interface{}) (*http.Response, error) {
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
