package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type getCipherResponse struct {
	Data string `json:"data"`
}

type getCipherRequest struct {
	Key string `json:"key"`
}

type GetCipher struct {
	Repository Repository
	Decrypter  Decrypter
}

func (g *GetCipher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType != contentTypeJSON {
		writeError(w, http.StatusUnsupportedMediaType, "unsupported Content-Type")
		return
	}

	w.Header().Set("Content-Type", contentTypeJSON)

	vars := mux.Vars(r)
	id := vars["id"]

	body := &getCipherRequest{}
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "decoding request body")
		return
	}

	cipher, err := g.Repository.FindByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	plainText, err := g.Decrypter.Decrypt(body.Key, cipher.CipherText)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "wrong key")
		return
	}

	cipherRes := &getCipherResponse{
		Data: plainText,
	}

	resBody, err := json.Marshal(cipherRes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "encoding response body")
		return
	}

	w.Write(resBody)
}
