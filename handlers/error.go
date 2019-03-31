package handlers

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Message string `json:"error"`
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)

	errRes := errorResponse{
		Message: msg,
	}

	json.NewEncoder(w).Encode(errRes)
}
