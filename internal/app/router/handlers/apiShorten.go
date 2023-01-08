// Package handlers provides HTTP handlers for the url-shortener.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Mldlr/url-shortener/internal/app/router/middleware"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/Mldlr/url-shortener/internal/app/utils/validators"
)

// APIShorten processes a request to shorten a URL and returns it in JSON format.
func APIShorten(repo storage.Repository, c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode the request body into a URL struct.
		var body *models.URL
		var err error
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "error reading request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		// Check if URL is valid.
		if !validators.IsURL(body.LongURL) {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}
		// Generate a short ID for the URL.
		body.ShortURL, err = repo.NewID(body.LongURL)
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting new id: %v", err), http.StatusInternalServerError)
			return
		}
		// Get the user ID from the request.
		userID, found := middleware.GetUserID(r)
		if !found {
			http.Error(w, fmt.Sprintf("error getting user cookie: %v", err), http.StatusInternalServerError)
			return
		}
		body.UserID = userID
		// Add the URL to the repository.
		duplicates, err := repo.Add(r.Context(), body)
		if err != nil {
			http.Error(w, fmt.Sprintf("error adding record to db: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if duplicates {
			// If the URL already exists, return a conflict status.
			w.WriteHeader(http.StatusConflict)
		} else {
			// If the URL is new, return a created status.
			w.WriteHeader(http.StatusCreated)
		}
		if err := json.NewEncoder(w).Encode(models.Response{Result: strings.Join([]string{c.BaseURL, body.ShortURL}, "/")}); err != nil {
			http.Error(w, "error building the response", http.StatusInternalServerError)
			return
		}
	}
}
