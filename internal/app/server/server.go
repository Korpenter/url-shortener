package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/Mldlr/url-shortener/internal/app/tls"
	"github.com/go-chi/chi/v5"
)

// Server is the structure to wrap http server config and shutdown channel for graceful shutdown.
type Server struct {
	cfg              *config.Config
	srv              http.Server
	shutdownFinished chan struct{}
}

// NewServer creates a new server instance
func NewServer(r chi.Router, c *config.Config) *Server {
	return &Server{
		cfg: c,
		srv: http.Server{
			Handler: r,
			Addr:    fmt.Sprintf(c.ServerAddress),
		},
		shutdownFinished: make(chan struct{}),
	}
}

// Run starts the server
func (s *Server) Run() {
	var err error
	if s.shutdownFinished == nil {
		s.shutdownFinished = make(chan struct{})
	}

	if s.cfg.EnableHTTPS {
		certFiles := []string{s.cfg.CertFile, s.cfg.KeyFile}
		for _, file := range certFiles {
			if _, err := os.Stat(file); err != nil {
				err = tls.GenerateCert(s.cfg)
				if err != nil {
					log.Fatalf("error generating certificate: %v", err)
				}
				break
			}
		}
		err = s.srv.ListenAndServeTLS(s.cfg.CertFile, s.cfg.KeyFile)
	} else {
		err = s.srv.ListenAndServe()
	}

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("unexpected error starting server %v", err)
	}

	log.Println("waiting for shutdown finishing...")
	<-s.shutdownFinished
	log.Println("shutdown finished")
}

// WaitForExitingSignal waits for a signal to exit the server and shutdowns it
func (s *Server) WaitForExitingSignal(timeout time.Duration, r storage.Repository) {
	var waiter = make(chan os.Signal, 1)
	signal.Notify(waiter, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-waiter

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := r.Close(); err != nil {
		log.Printf("failed to close repo: %v", err)
	}

	err := s.srv.Shutdown(ctx)
	if err != nil {
		log.Println("shutting down: " + err.Error())
	} else {
		log.Println("shutdown processed successfully")
		close(s.shutdownFinished)
	}
}
