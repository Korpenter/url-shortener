package encoders

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func HMACString(s string, k string) string {
	h := hmac.New(sha256.New, []byte(k))
	h.Write([]byte(s))
	hex := hex.EncodeToString(h.Sum(nil))
	fmt.Println("encoded: ", hex)
	return hex
}
