package main

import (
	"github.com/Mldlr/url-shortener/internal/app/shortener/server"
)

func main() {
	addr := "localhost:8080"
	server.StartServer(addr)
}
