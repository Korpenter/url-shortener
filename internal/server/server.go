package server

import (
	"github.com/Mldlr/url-shortener/internal/handler"
	"net/http"
)

const (
	Address = "localhost:8080"
	BaseURL = "http://localhost:8080"
)

func New() *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", handler.NewShortenerHandler())
	return &http.Server{
		Handler: mux,
		Addr:    Address,
	}
}
