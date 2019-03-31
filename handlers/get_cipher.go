package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type getCipherResponse struct {
	ResourceID string `json:"resource_id"`
	Data       string `json:"data"`
}

type getCipherRequest struct {
	Key string `json:"key"`
}

type GetCipher struct {
	Repository Repository
}

func (g *GetCipher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	resourceID := vars["resource_id"]

	body := &getCipherRequest{}
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "decoding request body")
		return
	}

	cipher, err := g.Repository.FindByResourceID(resourceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "finding cipher")
		return
	}

	if body.Key != cipher.Key {
		writeError(w, http.StatusUnauthorized, "wrong key")
		return
	}

	cipherRes := &getCipherResponse{
		ResourceID: cipher.ResourceID,
		Data:       cipher.CipherText,
	}

	resBody, err := json.Marshal(cipherRes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "encoding response body")
		return
	}

	w.Write(resBody)
}
