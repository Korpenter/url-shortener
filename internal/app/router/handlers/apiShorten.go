// Package handlers provides HTTP handlers for the url-shortener.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/service"
	"github.com/Mldlr/url-shortener/internal/app/utils/helpers"
)

// APIShorten processes a request to shorten a URL and returns it in JSON format.
func APIShorten(shortener service.ShortenerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, found := helpers.GetUserID(r)
		if !found {
			http.Error(w, "error getting user cookie", http.StatusInternalServerError)
			return
		}
		// Decode the request body into a URL struct.
		var body *models.URL
		var err error
		if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "error reading request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		var statusCode int
		body.UserID = userID
		url, err := shortener.Shorten(r.Context(), body)
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
			statusCode = http.StatusConflict
		} else {
			// If URL is new, return a created status.
			statusCode = http.StatusCreated
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if err := json.NewEncoder(w).Encode(models.Response{Result: shortener.BuildURL(url.ShortURL)}); err != nil {
			http.Error(w, "error building the response", http.StatusInternalServerError)
			return
		}
	}
}
