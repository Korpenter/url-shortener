package validators

import (
	"net/url"
	"strings"
)

func IsURL(longURL string) bool {
	if !strings.Contains(longURL, "://") {
		longURL = "http://" + longURL
	}
	_, err := url.ParseRequestURI(longURL)
	if err != nil {
		return false
	}
	u, err := url.Parse(longURL)
	if err != nil || u.Host == "" {
		return false
	}
	return true
}
