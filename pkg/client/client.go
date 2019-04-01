package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/jamesjoshuahill/ciphers/handlers"
)

const ciphersResourcePath = "/v1/ciphers"

type ServerClient interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
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

func (*client) Retrieve(id, aesKey []byte) ([]byte, error) {
	return nil, nil
}
