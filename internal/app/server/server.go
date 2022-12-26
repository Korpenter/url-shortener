package server

import (
	"fmt"
	"net/http"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/go-chi/chi/v5"
)

// NewServer initializes and HTTP server
func NewServer(r chi.Router, c *config.Config) *http.Server {
	return &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(c.ServerAddress),
	}
}
