package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type createCipherRequest struct {
	Data       string `json:"data"`
	ResourceID string `json:"resource_id"`
}

type createCipherResponse struct {
	ResourceID string `json:"resource_id"`
	Key        string `json:"key"`
}

func CreateCipherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body := &createCipherRequest{}
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, errorResponseBody(fmt.Sprintf("decoding request body: %s", err)))
		return
	}

	cipher := createCipherResponse{
		ResourceID: body.ResourceID,
		Key:        "key for server-cipher-id",
	}

	resBody, err := json.Marshal(cipher)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, errorResponseBody(fmt.Sprintf("encoding response body: %s", err)))
		return
	}

	w.Write(resBody)
}
