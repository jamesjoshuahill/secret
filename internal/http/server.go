package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	s *http.Server
}

func NewServer(host string, port int, h http.Handler) *Server {
	return &Server{
		s: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", host, port),
			Handler:      h,
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
		},
	}
}

func (s *Server) StartTLS(certFile, keyFile string) {
	started := make(chan struct{})

	go func() {
		started <- struct{}{}

		err := s.s.ListenAndServeTLS(certFile, keyFile)
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve error: %s", err)
		}
	}()

	<-started
}

func (s *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return s.s.Shutdown(ctx)
}
