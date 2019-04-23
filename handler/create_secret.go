package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jamesjoshuahill/secret/inmemory"
)

const contentTypeJSON = "application/json"

type CreateSecretRequest struct {
	Data string `json:"data"`
	ID   string `json:"id"`
}

type CreateSecretResponse struct {
	Key string `json:"key"`
}

type CreateSecret struct {
	Repository Repository
	Encrypt    EncryptFunc
}

func (c *CreateSecret) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType != contentTypeJSON {
		writeError(w, http.StatusUnsupportedMediaType, "unsupported Content-Type")
		return
	}

	w.Header().Set("Content-Type", contentTypeJSON)

	reqBody := &CreateSecretRequest{}
	err := json.NewDecoder(r.Body).Decode(reqBody)
	if err != nil {
		writeError(w, http.StatusBadRequest, "decoding request body")
		return
	}

	secret, err := c.Encrypt(reqBody.Data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "encrypting data")
		return
	}

	err = c.Repository.Store(inmemory.Secret{
		ID:         reqBody.ID,
		Nonce:      secret.Nonce,
		CipherText: secret.CipherText,
	})
	if err != nil {
		writeError(w, http.StatusConflict, "secret already exists")
		return
	}

	secretRes := CreateSecretResponse{
		Key: secret.Key,
	}

	err = json.NewEncoder(w).Encode(secretRes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "encoding response body")
		return
	}
}
