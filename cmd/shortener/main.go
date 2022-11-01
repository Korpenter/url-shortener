package main

import (
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/router"
	"github.com/Mldlr/url-shortener/internal/app/server"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"log"
)

func main() {
	cfg := config.NewConfig()
	repo := storage.New(cfg)
	r := router.NewRouter(repo, cfg)
	s := server.NewServer(r, cfg)
	log.Printf("Starting with cfg: %v", cfg)
	log.Fatal(s.ListenAndServe())
}
