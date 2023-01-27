package tls

import (
	"crypto/tls"
	"github.com/Mldlr/url-shortener/internal/app/config"
)

func NewTLSConfig(c *config.Config) (*tls.Config, error) {
	pubKey, privKey, err := Generate(options)
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(pubKey, privKey)
	if err != nil {
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{cert}}, nil
}
