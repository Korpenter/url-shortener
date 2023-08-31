package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/service"
)

// Expand redirects the client to the original URL associated with the short URL
// in the request path.
func Expand(shortener service.ShortenerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the short URL id from request path.
		id := strings.Split(r.URL.Path, "/")[1:]
		// Get the URL from the storage repository.
		url, err := shortener.Expand(r.Context(), id[0])
		if err != nil {
			// If the URL has been deleted, return Gone status.
			if !errors.Is(err, models.ErrURLDeleted) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusGone)
			return
		}
		// To redirect the client set the Location header to original URL.
		w.Header().Set("Location", url.LongURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
