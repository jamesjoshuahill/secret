package handler

import (
	"encoding/json"
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

	json.NewEncoder(w).Encode(errRes)
}
