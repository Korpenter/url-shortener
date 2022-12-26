package config

import (
	"flag"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerAddress string `envconfig:"SERVER_ADDRESS" default:"localhost:8080"`
	BaseURL       string `envconfig:"BASE_URL" default:"http://localhost:8080/"`
	FileStorage   string `envconfig:"FILE_STORAGE_PATH" default:""`
	SecretKey     string `envconfig:"URL_SHORTENER_KEY" default:"defaultKeyUrlSHoRtenEr"`
	PostgresURL   string `envconfig:"DATABASE_DSN" default:""`
}

// NewConfig returns a pointer to a new config instance.
func NewConfig() *Config {
	var c Config
	envconfig.MustProcess("", &c)
	flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "server address")
	flag.StringVar(&c.BaseURL, "b", c.BaseURL, "base url address")
	flag.StringVar(&c.FileStorage, "f", c.FileStorage, "storage path")
	flag.StringVar(&c.SecretKey, "k", c.SecretKey, "key")
	flag.StringVar(&c.PostgresURL, "d", c.PostgresURL, "postgres url")
	flag.Parse()
	return &c
}
