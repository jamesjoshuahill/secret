package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	port := flag.Int("port", 0, "Port to serve HTTP")
	flag.Parse()

	if *port == 0 {
		fmt.Println("--port flag required")
		os.Exit(1)
	}

	r := mux.NewRouter()
	r.HandleFunc("/v1/ciphers", createCipherHandler).Methods("POST")
	http.Handle("/", r)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      r,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("Starting server on port %d\n", *port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

type createCipherRequest struct {
	Data       string `json:"data"`
	ResourceID string `json:"resource_id"`
}

type createCipherResponse struct {
	ID         string `json:"id"`
	ResourceID string `json:"resource_id"`
	Key        string `json:"key"`
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
		ID:         "server cipher id",
		ResourceID: body.ResourceID,
		Key:        "key for server cipher id",
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
