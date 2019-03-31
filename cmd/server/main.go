package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")

	r := mux.NewRouter()
	r.HandleFunc("/v1/ciphers", createCipherHandler).Methods("POST")
	http.Handle("/", r)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	log.Fatalln(err)
}

type createCipherRequest struct {
	Data       string `json:"data"`
	ResourceID string `json:"resource_id"`
}

type createCipherResponse struct {
	ID string `json:"id"`
	createCipherRequest
}

type errorResponse struct {
	Message string `json:"error"`
}

func createCipherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body := &createCipherRequest{}
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, errorResponseBody("decoding request body", err))
		return
	}

	cipher := createCipherResponse{
		ID:                  "server cipher id",
		createCipherRequest: *body,
	}

	resBody, err := json.Marshal(cipher)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, errorResponseBody("encoding response body", err))
		return
	}

	w.Write(resBody)
}

func errorResponseBody(context string, err error) string {
	errRes := errorResponse{
		Message: fmt.Sprintf("%s: %s", context, err),
	}

	resBody, err := json.Marshal(errRes)
	if err != nil {
		return fmt.Sprintf(`{"error":"encoding error response body: %s"}`, err)
	}

	return string(resBody)
}
