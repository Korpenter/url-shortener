package main

import (
	"fmt"
	"log"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/router"
	"github.com/Mldlr/url-shortener/internal/app/server"
	"github.com/Mldlr/url-shortener/internal/app/storage"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

const NA string = "N/A"

func main() {
	if len(buildVersion) == 0 {
		buildVersion = NA
	}
	if len(buildDate) == 0 {
		buildDate = NA
	}
	if len(buildCommit) == 0 {
		buildCommit = NA
	}
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	cfg := config.NewConfig()
	repo := storage.New(cfg)
	r := router.NewRouter(repo, cfg)
	s := server.NewServer(r, cfg)
	log.Printf("Starting with cfg: %v", cfg)
	log.Fatal(s.ListenAndServe())
}
