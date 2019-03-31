package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jamesjoshuahill/ciphers/repository"
)

type createCipherRequest struct {
	Data       string `json:"data"`
	ResourceID string `json:"resource_id"`
}

type createCipherResponse struct {
	ResourceID string `json:"resource_id"`
	Key        string `json:"key"`
}

type CreateCipher struct {
	Repository Repository
}

func (c *CreateCipher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reqBody := &createCipherRequest{}
	err := json.NewDecoder(r.Body).Decode(reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, errorResponseBody(fmt.Sprintf("decoding request body: %s", err)))
		return
	}

	key := "key for server-cipher-id"

	_ = c.Repository.Store(repository.Cipher{
		ResourceID: reqBody.ResourceID,
		Data:       reqBody.Data,
		Key:        key,
	})

	cipher := createCipherResponse{
		ResourceID: reqBody.ResourceID,
		Key:        key,
	}

	resBody, err := json.Marshal(cipher)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, errorResponseBody(fmt.Sprintf("encoding response body: %s", err)))
		return
	}

	w.Write(resBody)
}
