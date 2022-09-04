package main

import (
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/server"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"log"
)

func main() {
	cfg := config.New()
	repo := storage.NewInMemRepo()
	s := server.New(repo, cfg)
	log.Fatal(s.ListenAndServe())
}
