package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jamesjoshuahill/secret/aes"
	"github.com/jamesjoshuahill/secret/handler"
	"github.com/jamesjoshuahill/secret/inmemory"

	"github.com/gorilla/mux"
	"github.com/jessevdk/go-flags"
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

	repo := inmemory.NewRepo()
	createSecretHandler := &handler.CreateSecret{Repository: repo, Encrypt: aes.Encrypt}
	getSecretHandler := &handler.GetSecret{Repository: repo, Decrypt: aes.Decrypt}

	r := mux.NewRouter()
	r.Methods(http.MethodPost).Path("/v1/secrets").Handler(createSecretHandler)
	r.Methods(http.MethodGet).Path("/v1/secrets/{id}").Handler(getSecretHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", opts.Port),
		Handler:      r,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	serverErr := make(chan error, 1)
	log.Printf("starting server on port %d\n", opts.Port)
	go func() {
		err := server.ListenAndServeTLS(opts.Cert, opts.Key)
		if err != nil && err != http.ErrServerClosed {
			serverErr <- fmt.Errorf("listen and serve error: %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-stop:
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if err := server.Shutdown(ctx); err != nil {
				serverErr <- fmt.Errorf("shutdown error: %s", err)
			} else {
				log.Println("shut down server gracefully")
				os.Exit(0)
			}
		case err := <-serverErr:
			log.Fatalf("%s\n", err)
		}
	}
}
