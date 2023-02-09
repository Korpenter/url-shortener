package tls

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/Mldlr/url-shortener/internal/app/config"
)

// GenerateCert generates a new pair of tls certificate and key and writes them to .pem files.
func GenerateCert(c *config.Config) error {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		err = fmt.Errorf("failed to generate private key: %v", err)
		return err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		err = fmt.Errorf("failed to generate serial number: %v", err)
		return err
	}

	keyBytes := x509.MarshalPKCS1PublicKey(&priv.PublicKey)
	keyHash := sha1.Sum(keyBytes)
	ski := keyHash[:]

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"shortener"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: ski,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, priv.Public(), priv)
	if err != nil {
		err = fmt.Errorf("failed to create certificate: %v", err)
		return err
	}

	var certBuf bytes.Buffer
	err = pem.Encode(&certBuf, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	if err != nil {
		err = fmt.Errorf("failed to encode cert to pem: %v", err)
		return err
	}

	err = os.WriteFile(c.CertFile, certBuf.Bytes(), 0644)
	if err != nil {
		err = fmt.Errorf("failed to write cert key to file: %v", err)
		return err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		err = fmt.Errorf("failed to marshal private key: %v", err)
		return err
	}

	var privKeyBuf bytes.Buffer
	err = pem.Encode(&privKeyBuf, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	if err != nil {
		err = fmt.Errorf("failed to encode private key to pem: %v", err)
		return err
	}

	err = os.WriteFile(c.KeyFile, privKeyBuf.Bytes(), 0644)
	if err != nil {
		err = fmt.Errorf("failed to write private key to file: %v", err)
		return err
	}

	return nil
}
