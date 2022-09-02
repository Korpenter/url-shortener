package utils

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
)

func MD5(data string) string {
	h := md5.Sum([]byte(data))
	return hex.EncodeToString(h[:])
}

func IsValid(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}
