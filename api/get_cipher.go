package api

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

func (*api) GetCipherHandleFunc(w http.ResponseWriter, r *http.Request) {
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

	if body.Key != "key for server-cipher-id" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, errorResponseBody("wrong key"))
		return
	}

	cipher := &getCipherResponse{
		ResourceID: resourceID,
		Data:       "some plain text",
	}

	resBody, err := json.Marshal(cipher)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, errorResponseBody(fmt.Sprintf("encoding response body: %s", err)))
		return
	}

	w.Write(resBody)
}
