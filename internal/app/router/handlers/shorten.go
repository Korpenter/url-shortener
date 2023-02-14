package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/service"
	"github.com/Mldlr/url-shortener/internal/app/utils/helpers"
)

// Shorten creates a new shortened URL.
func Shorten(shortener service.ShortenerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from request.
		userID, found := helpers.GetUserID(r)
		if !found {
			http.Error(w, "error getting user cookie", http.StatusInternalServerError)
			return
		}
		// Read request body.
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request: %s", err.Error()), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		var status int
		url, err := shortener.Shorten(r.Context(), &models.URL{LongURL: string(b), UserID: userID})
		// Create URL model, and add it to storage.
		if err != nil {
			// If there is an error, and its not a duplicate url
			if errors.Is(err, models.ErrInvalidURL) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			} else if !errors.Is(err, models.ErrDuplicate) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// If it is a duplicate url
			status = http.StatusConflict
		} else {
			// If URL is new, return a created status.
			status = http.StatusCreated
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(status)
		if _, err = io.WriteString(w, shortener.BuildURL(url.ShortURL)); err != nil {
			http.Error(w, "error building the response", http.StatusInternalServerError)
		}
	}
}
