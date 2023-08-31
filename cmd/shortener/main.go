package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Mldlr/url-shortener/internal/app/config"
	grpc "github.com/Mldlr/url-shortener/internal/app/grpc/server"
	"github.com/Mldlr/url-shortener/internal/app/router"
	"github.com/Mldlr/url-shortener/internal/app/server"
	"github.com/Mldlr/url-shortener/internal/app/service"
	"github.com/Mldlr/url-shortener/internal/app/storage"
)

// Build info
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// NA is the string output if build info is not set
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
	shortener := service.NewShortenerImpl(repo, cfg)
	r := router.NewRouter(shortener, cfg)
	s := server.NewServer(r, cfg)
	log.Printf("Starting with cfg: %v", cfg)
	if cfg.GRPCAddress != "" {
		grpcS := grpc.NewGRPCServer(shortener, cfg)
		go grpcS.Run(context.Background())
	}
	go s.WaitForExitingSignal(15*time.Second, repo)
	s.Run()
}
