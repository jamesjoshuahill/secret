package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jamesjoshuahill/ciphers/handlers"
)

// Retrieve sends an HTTP request to get the cipher using and id and aesKey,
// and returns the decrypted plain text.
func (c *Client) Retrieve(id, aesKey []byte) ([]byte, error) {
	reqBody := handlers.GetCipherRequest{
		Key: string(aesKey),
	}

	url := fmt.Sprintf("%s%s/%s", c.baseURL, ciphersResourcePath, string(id))
	res, err := c.do("GET", url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("get cipher request: %s", err)
	}

	switch res.StatusCode {
	case http.StatusOK:
		var body handlers.GetCipherResponse
		err = json.NewDecoder(res.Body).Decode(&body)
		if err != nil {
			return nil, fmt.Errorf("decoding get cipher response body: %s", err)
		}

		return []byte(body.Data), nil
	case http.StatusUnauthorized:
		return nil, wrongKeyError{}
	case http.StatusNotFound:
		return nil, notFoundError{}
	default:
		return nil, newUnexpectedError(res)
	}
}
