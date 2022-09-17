package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	ServerAddress string `envconfig:"SERVER_ADDRESS" default:"0.0.0.0:8080"`
	BaseUrl       string `envconfig:"BASE_URL" default:"http://localhost:8080/"`
}

// NewConfig returns a pointer to a new config instance.
func NewConfig() *Config {
	var c Config
	envconfig.MustProcess("", &c)
	return &c
}
