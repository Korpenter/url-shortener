package server

import (
	"github.com/Mldlr/url-shortener/internal/config"
	"github.com/Mldlr/url-shortener/internal/handler"
	"net/http"
)

func New() *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", handler.NewShortenerHandler())
	return &http.Server{
		Handler: mux,
		Addr:    config.Address,
	}
}
