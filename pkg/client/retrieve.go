package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jamesjoshuahill/secret/handler"
)

// Retrieve sends an HTTP request to get the secret using and id and aesKey,
// and returns the decrypted plain text.
func (c *Client) Retrieve(id, aesKey []byte) ([]byte, error) {
	reqBody := handler.GetSecretRequest{
		Key: string(aesKey),
	}

	url := fmt.Sprintf("%s%s/%s", c.baseURL, secretsResourcePath, string(id))
	res, err := c.do("GET", url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("get secret request: %s", err)
	}

	switch res.StatusCode {
	case http.StatusOK:
		var body handler.GetSecretResponse
		defer res.Body.Close() //nolint:errcheck
		err = json.NewDecoder(res.Body).Decode(&body)
		if err != nil {
			return nil, fmt.Errorf("decoding get secret response body: %s", err)
		}

		return []byte(body.Data), nil
	case http.StatusUnprocessableEntity:
		return nil, ErrWrongIDOrKey
	default:
		return nil, newUnexpectedError(res)
	}
}
