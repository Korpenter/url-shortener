package server

import (
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/handlers"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"net/http"
)

func New(r storage.Repositories, c *config.Config) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", handlers.NewShortenerHandler(r, c))
	return &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf("%s:%d", c.Host, c.Port),
	}
}
