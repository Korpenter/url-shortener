package server

import (
	"github.com/Mldlr/url-shortener/internal/app/shortener/handler"
	"log"
	"net/http"
)

type urlShortenerServer struct {
	server http.Server
}

func newShortenerServer(addr string) *urlShortenerServer {
	mux := http.NewServeMux()
	mux.Handle("/", handler.NewShortenerHandler())
	return &urlShortenerServer{
		server: http.Server{
			Handler: mux,
			Addr:    addr,
		},
	}
}

func StartServer(Addr string) {
	s := newShortenerServer(Addr)
	log.Fatal(s.server.ListenAndServe())
}
