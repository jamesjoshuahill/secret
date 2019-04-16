// Package client provides a Client for secret-server.
package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const secretsResourcePath = "/v1/secrets"

// Client represents an API client that provides store and retreive functions.
type Client struct {
	baseURL     string
	httpsClient HTTPSClient
}

// New returns a Client struct with the given baseURL of the secret-server and HTTPSClient.
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
