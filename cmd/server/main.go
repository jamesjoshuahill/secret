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
	r.HandleFunc("/v1/ciphers/{resource_id}", getCipherHandler).Methods("GET")

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

type getCipherResponse struct {
	ResourceID string `json:"resource_id"`
	Data       string `json:"data"`
}

type getCipherRequest struct {
	Key string `json:"key"`
}

func getCipherHandler(w http.ResponseWriter, r *http.Request) {
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
