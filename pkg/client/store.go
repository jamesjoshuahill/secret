package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jamesjoshuahill/ciphers/handler"
)

// Store sends an HTTP request to create a cipher of the payload with an id,
// and returns the aesKey.
func (c *Client) Store(id, payload []byte) ([]byte, error) {
	reqBody := handler.CreateSecretRequest{
		ID:   string(id),
		Data: string(payload),
	}

	res, err := c.do("POST", c.baseURL+ciphersResourcePath, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create cipher request: %s", err)
	}

	switch res.StatusCode {
	case http.StatusOK:
		var body handler.CreateSecretResponse
		err = json.NewDecoder(res.Body).Decode(&body)
		if err != nil {
			return nil, fmt.Errorf("decoding create cipher response body: %s", err)
		}

		return []byte(body.Key), nil
	case http.StatusConflict:
		return nil, alreadyExistsError{}
	default:
		return nil, newUnexpectedError(res)
	}
}
