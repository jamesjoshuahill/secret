package handlers

import (
	"encoding/json"
	"fmt"
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
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, errorResponseBody(fmt.Sprintf("decoding request body: %s", err)))
		return
	}

	cipher, _ := g.Repository.FindByResourceID(resourceID)

	if body.Key != cipher.Key {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, errorResponseBody("wrong key"))
		return
	}

	cipherRes := &getCipherResponse{
		ResourceID: cipher.ResourceID,
		Data:       cipher.Data,
	}

	resBody, err := json.Marshal(cipherRes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, errorResponseBody(fmt.Sprintf("encoding response body: %s", err)))
		return
	}

	w.Write(resBody)
}
