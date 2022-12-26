// Package validators provides functions to validate data.
package validators

import (
	"net/url"
	"strings"
)

// IsURL checks if url is a valid and defaults it to http if method is not presented.
func IsURL(longURL string) bool {
	// Check if url contains proto.
	if !strings.Contains(longURL, "://") {
		// Default it to http if it doesnt.
		longURL = "http://" + longURL
	}
	// Check url with stdlib function.
	u, err := url.Parse(longURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}
