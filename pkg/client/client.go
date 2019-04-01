package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jamesjoshuahill/ciphers/handlers"
)

const ciphersResourcePath = "/v1/ciphers"

type ServerClient interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
	Do(req *http.Request) (*http.Response, error)
}

type client struct {
	baseURL      string
	serverClient ServerClient
}

func New(baseURL string, serverClient ServerClient) *client {
	return &client{
		baseURL:      baseURL,
		serverClient: serverClient,
	}
}

func (c *client) Store(id, payload []byte) ([]byte, error) {
	reqBody := handlers.CreateCipherRequest{
		ID:   string(id),
		Data: string(payload),
	}

	reqBytes, _ := json.Marshal(&reqBody)

	res, _ := c.serverClient.Post(c.baseURL+ciphersResourcePath, "application/json", bytes.NewReader(reqBytes))

	var body handlers.CreateCipherResponse
	_ = json.NewDecoder(res.Body).Decode(&body)

	return []byte(body.Key), nil
}

func (c *client) Retrieve(id, aesKey []byte) ([]byte, error) {
	reqBody := handlers.GetCipherRequest{
		Key: string(aesKey),
	}

	reqBytes, _ := json.Marshal(&reqBody)

	url := fmt.Sprintf("%s%s/%s", c.baseURL, ciphersResourcePath, string(id))
	req, _ := http.NewRequest("GET", url, bytes.NewReader(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	res, _ := c.serverClient.Do(req)

	var body handlers.GetCipherResponse
	_ = json.NewDecoder(res.Body).Decode(&body)

	return []byte(body.Data), nil
}
