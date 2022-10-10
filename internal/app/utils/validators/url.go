package validators

import (
	"net/url"
	"strings"
)

// IsURL checks if url is a valid and default it to http if method is not presented.
func IsURL(longURL string) bool {
	if !strings.Contains(longURL, "://") {
		longURL = "http://" + longURL
	}
	u, err := url.Parse(longURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}
