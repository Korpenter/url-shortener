package server

import (
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func New(r chi.Router, c *config.Config) *http.Server {
	return &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s:%d", c.Host, c.Port),
	}
}
