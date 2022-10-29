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
	sha.Write([]byte(url))
	num := new(big.Int).SetBytes(sha.Sum(nil)).Uint64()
	var b strings.Builder
	for num > 0 {
		r := num % base
		num /= base
		b.WriteString(string(charSet[r]))
	}
	return b.String()
}
