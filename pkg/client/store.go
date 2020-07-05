package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jamesjoshuahill/secret/handler"
)

// Store sends an HTTP request to create a secret of the payload with an id,
// and returns the aesKey.
func (c *Client) Store(id, payload []byte) ([]byte, error) {
	reqBody := handler.CreateSecretRequest{
		ID:   string(id),
		Data: string(payload),
	}

	res, err := c.do("POST", c.baseURL+secretsResourcePath, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create secret request: %s", err)
	}

	switch res.StatusCode {
	case http.StatusOK:
		var body handler.CreateSecretResponse
		defer res.Body.Close() //nolint:errcheck
		err = json.NewDecoder(res.Body).Decode(&body)
		if err != nil {
			return nil, fmt.Errorf("decoding create secret response body: %s", err)
		}

		return []byte(body.Key), nil
	case http.StatusConflict:
		return nil, ErrAlreadyExists
	default:
		return nil, newUnexpectedError(res)
	}
}
