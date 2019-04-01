package client

import (
	"encoding/json"
	"fmt"

	"github.com/jamesjoshuahill/ciphers/handlers"
)

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
