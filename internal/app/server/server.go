package server

import (
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// NewServer returns a pointer to a new http.Server instance
func NewServer(r chi.Router, c *config.Config) *http.Server {
	return &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(c.ServerAddress),
	}
}
