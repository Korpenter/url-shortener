package encode

import (
	"strings"
)

const (
	base    = 62
	charSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

func ToBase62(num int) string {
	var b strings.Builder
	for num > 0 {
		r := num % base
		num /= base
		b.WriteString(string(charSet[r]))
	}
	return b.String()
}
