// Package encoders provides functions to encode data.
package encoders

import (
	"crypto/sha256"
	"math/big"
	"strings"
)

const (
	base    = 62
	charSet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// ToRBase62 encodes a number to Base62 string in reversed order
func ToRBase62(url string) string {
	sha := sha256.New()
	// Calculate hash.
	sha.Write([]byte(url))
	num := new(big.Int).SetBytes(sha.Sum(nil)).Uint64()
	// Build string ID.
	var b strings.Builder
	b.Grow(64)
	for num > 0 {
		r := num % base
		num /= base
		b.WriteString(string(charSet[r]))
	}
	return b.String()
}
