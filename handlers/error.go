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
		return fmt.Sprintf(`{"error":"encoding error response body: %s"}`, err)
	}

	return string(body)
}
