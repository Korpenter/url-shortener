// Package config provides a configuration structure for the server and its initialization.
package config

import (
	"encoding/json"
	"flag"
	"log"
	"net/netip"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// Config represents the configuration options for the server.
type Config struct {
	ServerAddress string `envconfig:"SERVER_ADDRESS" default:"localhost:8080" json:"server_address"`
	BaseURL       string `envconfig:"BASE_URL" default:"http://localhost:8080" json:"base_url"`
	FileStorage   string `envconfig:"FILE_STORAGE_PATH" default:"" json:"file_storage_path"`
	SecretKey     []byte `envconfig:"URL_SHORTENER_KEY" default:"defaultKeyUrlSHoRtenEr" json:"secret_key"`
	PostgresURL   string `envconfig:"DATABASE_DSN" default:"" json:"database_dsn"`
	EnableHTTPS   bool   `envconfig:"ENABLE_HTTPS" default:"false" json:"enable_https"`
	CertFile      string `envconfig:"TLS_CERT_FILE" default:"cert.pem" json:"cert_file"`
	KeyFile       string `envconfig:"TLS_KEY_FILE" default:"key.pem" json:"key_file"`
	TrustedSubnet string `envconfig:"TRUSTED_SUBNET" default:"" json:"trusted_subnet"`
	SubnetPrefix  netip.Prefix
}

// NewConfig initializes and returns a new Config struct. It reads
// environment variables and command-line flags to set the configuration values.
func NewConfig() *Config {
	var c Config
	var err error
	configFile := os.Getenv("CONFIG")
	envconfig.MustProcess("", &c)

	flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "server address")
	flag.StringVar(&c.BaseURL, "b", c.BaseURL, "base url address")
	flag.StringVar(&c.FileStorage, "f", c.FileStorage, "storage path")
	flag.StringVar(&c.PostgresURL, "d", c.PostgresURL, "postgres url")
	flag.BoolVar(&c.EnableHTTPS, "s", c.EnableHTTPS, "enable https")
	flag.StringVar(&c.CertFile, "l", c.CertFile, "tls cert file path")
	flag.StringVar(&c.KeyFile, "key", c.KeyFile, "tls key file path")
	flag.StringVar(&configFile, "c", configFile, "path to config file")
	flag.StringVar(&c.TrustedSubnet, "t", c.TrustedSubnet, "trusted subnet CIDR")
	key := flag.String("k", "", "key")
	flag.Parse()

	if configFile != "" {
		bytes, err := os.ReadFile(configFile)
		if err != nil {
			log.Fatalf("failed to read cfg from file: %v", err)
		}

		err = json.Unmarshal(bytes, &c)
		if err != nil {
			log.Fatalf("failed to marshal cfg from file: %v", err)
		}
	}

	if c.EnableHTTPS {
		c.BaseURL = strings.Replace(c.BaseURL, "http://", "https://", 1)
	}

	if c.TrustedSubnet != "" {
		c.SubnetPrefix, err = netip.ParsePrefix(c.TrustedSubnet)
		if err != nil {
			log.Fatalf("invalid trusted network: %v", err)
		}
	}
	if *key != "" {
		c.SecretKey = []byte(*key)
	}
	return &c
}
