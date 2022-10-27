package encoders

import (
	"strings"
)

const (
	base    = 62
	charSet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// ToRBase62 encodes a number to Base62 string in reversed order
func ToRBase62(num int) string {
	var b strings.Builder
	for num > 0 {
		r := num % base
		num /= base
		b.WriteString(string(charSet[r]))
	}
	return b.String()
}
