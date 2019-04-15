package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jamesjoshuahill/ciphers/repository"
)

const contentTypeJSON = "application/json"

type CreateCipherRequest struct {
	Data string `json:"data"`
	ID   string `json:"id"`
}

type CreateCipherResponse struct {
	Key string `json:"key"`
}

type CreateCipher struct {
	Repository Repository
	Encrypter  Encrypter
}

func (c *CreateCipher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType != contentTypeJSON {
		writeError(w, http.StatusUnsupportedMediaType, "unsupported Content-Type")
		return
	}

	w.Header().Set("Content-Type", contentTypeJSON)

	reqBody := &CreateCipherRequest{}
	err := json.NewDecoder(r.Body).Decode(reqBody)
	if err != nil {
		writeError(w, http.StatusBadRequest, "decoding request body")
		return
	}

	cipher, err := c.Encrypter.Encrypt(reqBody.Data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "encrypting data")
		return
	}

	err = c.Repository.Store(repository.Cipher{
		ID:         reqBody.ID,
		Nonce:      cipher.Nonce,
		CipherText: cipher.CipherText,
	})
	if err != nil {
		writeError(w, http.StatusConflict, "cipher already exists")
		return
	}

	cipherRes := CreateCipherResponse{
		Key: cipher.Key,
	}

	resBody, err := json.Marshal(cipherRes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "encoding response body")
		return
	}

	w.Write(resBody)
}
