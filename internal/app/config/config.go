// Package config provides a configuration structure for the server and its initialization.
package config

import (
	"flag"

	"github.com/kelseyhightower/envconfig"
)

// Config represents the configuration options for the server.
type Config struct {
	ServerAddress string `envconfig:"SERVER_ADDRESS" default:"localhost:8080"`
	BaseURL       string `envconfig:"BASE_URL" default:"http://localhost:8080"`
	FileStorage   string `envconfig:"FILE_STORAGE_PATH" default:""`
	SecretKey     []byte `envconfig:"URL_SHORTENER_KEY" default:"defaultKeyUrlSHoRtenEr"`
	PostgresURL   string `envconfig:"DATABASE_DSN" default:""`
}

// NewConfig initializes and returns a new Config struct. It reads
// environment variables and command-line flags to set the configuration values.
func NewConfig() *Config {
	var c Config
	envconfig.MustProcess("", &c)
	flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "server address")
	flag.StringVar(&c.BaseURL, "b", c.BaseURL, "base url address")
	flag.StringVar(&c.FileStorage, "f", c.FileStorage, "storage path")
	flag.StringVar(&c.PostgresURL, "d", c.PostgresURL, "postgres url")
	key := flag.String("k", "", "key")
	if *key != "" {
		c.SecretKey = []byte(*key)
	}
	flag.Parse()
	return &c
}
