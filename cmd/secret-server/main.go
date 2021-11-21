package main

import (
	"log"
	"os"
	"syscall"
	"time"

	"github.com/jamesjoshuahill/secret/internal/aes"
	"github.com/jamesjoshuahill/secret/internal/handler"
	"github.com/jamesjoshuahill/secret/internal/http"
	"github.com/jamesjoshuahill/secret/internal/inmemory"
	"github.com/jamesjoshuahill/secret/internal/signal"

	"github.com/gorilla/mux"
	"github.com/jessevdk/go-flags"
)

type options struct {
	Port int    `long:"port" description:"Port to serve HTTPS" required:"true"`
	Cert string `long:"cert" description:"Path to TLS certificate file" required:"true"`
	Key  string `long:"key" description:"Path to TLS private key file" required:"true"`
}

func main() {
	opts := parseOptions()

	repo := inmemory.NewRepo()
	createSecretHandler := &handler.CreateSecret{Repository: repo, Encrypt: aes.Encrypt}
	getSecretHandler := &handler.GetSecret{Repository: repo, Decrypt: aes.Decrypt}

	r := mux.NewRouter()
	r.Methods("POST").Path("/v1/secrets").Handler(createSecretHandler)
	r.Methods("GET").Path("/v1/secrets/{id}").Handler(getSecretHandler)

	server := http.NewServer(opts.Port, r)

	log.Printf("starting server on port %d\n", opts.Port)
	server.StartTLS(opts.Cert, opts.Key)

	signal.Wait(os.Interrupt, syscall.SIGTERM)
	log.Println("shutting down server...")

	if err := server.Shutdown(6 * time.Second); err != nil {
		log.Fatalf("error shutting down server: %s", err)
	}
}

func parseOptions() *options {
	opts := options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		if outErr, ok := err.(*flags.Error); ok && outErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	return &opts
}
