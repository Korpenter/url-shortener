package validate

import (
	"net/url"
	"strings"
)

func IsURL(LongURL string) bool {
	if !strings.Contains(LongURL, "://") {
		LongURL = "http://" + LongURL
	}
	_, err := url.ParseRequestURI(LongURL)
	if err != nil {
		return false
	}
	u, err := url.Parse(LongURL)
	if err != nil || u.Host != "" {
		return false
	}
	return true
}
