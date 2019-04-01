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

	switch res.StatusCode {
	case http.StatusOK:
		var body handlers.CreateCipherResponse
		err = json.NewDecoder(res.Body).Decode(&body)
		if err != nil {
			return nil, fmt.Errorf("decoding create cipher response body: %s", err)
		}

		return []byte(body.Key), nil
	case http.StatusConflict:
		return nil, alreadyExistsError{}
	default:
		unerr := unexpectedError{statusCode: res.StatusCode}

		var body handlers.ErrorResponse
		err = json.NewDecoder(res.Body).Decode(&body)
		if err != nil {
			return nil, unerr
		}

		unerr.message = body.Message
		return nil, unerr
	}
}

func (c *client) Retrieve(id, aesKey []byte) ([]byte, error) {
	reqBody := handlers.GetCipherRequest{
		Key: string(aesKey),
	}

	url := fmt.Sprintf("%s%s/%s", c.baseURL, ciphersResourcePath, string(id))
	res, err := c.do("GET", url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("get cipher request: %s", err)
	}

	var body handlers.GetCipherResponse
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("decoding get cipher response body: %s", err)
	}

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
