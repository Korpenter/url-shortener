package main

import (
	"fmt"
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
	fmt.Println(repo)
	s := server.NewServer(r, cfg)
	log.Fatal(s.ListenAndServe())
}
