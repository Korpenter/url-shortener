package main

import (
	"github.com/Mldlr/url-shortener/internal/server"
	"log"
)

func main() {
	s := server.New()
	log.Fatal(s.ListenAndServe())
}
