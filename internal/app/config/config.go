package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port   int    `envconfig:"PORT" default:"8080"`
	Host   string `envconfig:"HOST" default:"0.0.0.0"`
	Prefix string `envconfig:"PREFIX" default:"http://localhost:8080/"`
}

// NewConfig returns a pointer to a new config instance.
func NewConfig() *Config {
	var c Config
	envconfig.MustProcess("", &c)
	return &c
}
