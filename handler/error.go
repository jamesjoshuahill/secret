package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"error"`
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)

	errRes := ErrorResponse{
		Message: msg,
	}

	err := json.NewEncoder(w).Encode(errRes)
	if err != nil {
		log.Printf("encoding error response: %s", err)
	}
}
