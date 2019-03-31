package main

import (
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

func createCipherHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
