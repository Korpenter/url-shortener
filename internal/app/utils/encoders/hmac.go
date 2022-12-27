// Package encoders provides functions to encode data.
package encoders

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// HMACString calculates a HMAC-SHA256 hash of the input string, using the given key.
func HMACString(s string, k []byte) string {
	h := hmac.New(sha256.New, k)
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
