// Package encoders provides functions to encode data.
package encoders

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HMACString(s string, k string) string {
	h := hmac.New(sha256.New, []byte(k))
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
