package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jamesjoshuahill/ciphers/handlers"
)

const ciphersResourcePath = "/v1/ciphers"

type client struct {
	baseURL     string
	httpsClient HTTPSClient
}

func New(baseURL string, httpsClient HTTPSClient) *client {
	return &client{
		baseURL:     baseURL,
		httpsClient: httpsClient,
	}
}

func (c *client) Store(id, payload []byte) ([]byte, error) {
	reqBody := handlers.CreateCipherRequest{
		ID:   string(id),
		Data: string(payload),
	}

	res, err := c.do("POST", c.baseURL+ciphersResourcePath, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create cipher request: %s", err)
	}

	var body handlers.CreateCipherResponse
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("decoding create cipher response body: %s", err)
	}

	return []byte(body.Key), nil
}

func (c *client) Retrieve(id, aesKey []byte) ([]byte, error) {
	reqBody := handlers.GetCipherRequest{
		Key: string(aesKey),
	}

	url := fmt.Sprintf("%s%s/%s", c.baseURL, ciphersResourcePath, string(id))
	res, _ := c.do("GET", url, reqBody)

	var body handlers.GetCipherResponse
	_ = json.NewDecoder(res.Body).Decode(&body)

	return []byte(body.Data), nil
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
