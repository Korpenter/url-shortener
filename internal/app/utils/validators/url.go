package validators

import (
	"net/url"
	"strings"
)

func IsURL(longURL string) bool {
	if !strings.Contains(longURL, "://") {
		longURL = "http://" + longURL
	}
	u, err := url.Parse(longURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}
