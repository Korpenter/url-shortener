// Package handlers provides HTTP handlers for the url-shortener.
package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/router/middleware"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/Mldlr/url-shortener/internal/app/utils/validators"
)

// Shorten creates a new shortened URL.
func Shorten(repo storage.Repository, c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read request body.
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error reading request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		// Check if body is a valid URL.
		long := string(b)
		if !validators.IsURL(long) {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}
		// Get new short URL ID.
		id, err := repo.NewID(long)
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting new id: %v", err), http.StatusInternalServerError)
			return
		}
		// Get user ID from request.
		userID, found := middleware.GetUserID(r)
		if !found {
			http.Error(w, fmt.Sprintf("error getting user cookie: %v", err), http.StatusInternalServerError)
			return
		}
		// Create URL model, and add it to storage.
		url := models.URL{ShortURL: id, LongURL: long, UserID: userID, Deleted: false}
		duplicates, err := repo.Add(r.Context(), &url)
		if err != nil {
			http.Error(w, fmt.Sprintf("error adding record to db: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		if duplicates {
			// If any of the URLs already exist, return a conflict status.
			w.WriteHeader(http.StatusConflict)
		} else {
			// If all URLs are new, return a created status.
			w.WriteHeader(http.StatusCreated)
		}
		if _, err = io.WriteString(w, strings.Join([]string{c.BaseURL, url.ShortURL}, "/")); err != nil {
			log.Println(err)
		}
	}
}
