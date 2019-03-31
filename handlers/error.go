package handlers

import (
	"encoding/json"
	"fmt"
)

type errorResponse struct {
	Message string `json:"error"`
}

func errorResponseBody(msg string) string {
	errRes := errorResponse{
		Message: msg,
	}

	body, err := json.Marshal(errRes)
	if err != nil {
		return fmt.Sprintf(`{"error":"%s"}`, msg)
	}

	return string(body)
}
