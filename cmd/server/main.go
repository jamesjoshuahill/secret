package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jamesjoshuahill/ciphers/repository/inmemory"

	"github.com/jamesjoshuahill/ciphers/encryption"

	"github.com/jamesjoshuahill/ciphers/handler"

	"github.com/jessevdk/go-flags"

	"github.com/gorilla/mux"
)

type options struct {
	Port int    `long:"port" description:"Port to serve HTTPS" required:"true"`
	Cert string `long:"cert" description:"Path to TLS certificate file" required:"true"`
	Key  string `long:"key" description:"Path to TLS private key file" required:"true"`
}

func main() {
	opts := &options{}
	_, err := flags.Parse(opts)
	if err != nil {
		if outErr, ok := err.(*flags.Error); ok && outErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	repo := inmemory.New()
	encrypter := encryption.Encrypter{}
	decrypter := encryption.Decrypter{}
	createSecretHandler := &handler.CreateSecret{Repository: repo, Encrypter: encrypter}
	getSecretHandler := &handler.GetSecret{Repository: repo, Decrypter: decrypter}

	r := mux.NewRouter()
	r.Methods("POST").Path("/v1/ciphers").Handler(createSecretHandler)
	r.Methods("GET").Path("/v1/ciphers/{id}").Handler(getSecretHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", opts.Port),
		Handler:      r,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("Starting server on port %d\n", opts.Port)
	err = srv.ListenAndServeTLS(opts.Cert, opts.Key)
	if err != nil {
		log.Fatalln(err)
	}
}
