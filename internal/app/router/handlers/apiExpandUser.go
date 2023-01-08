package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/router/middleware"
	"github.com/Mldlr/url-shortener/internal/app/storage"
)

// APIUserExpand retrieves the list of shortened URLs
// created by a user and returns them as a JSON array.
func APIUserExpand(repo storage.Repository, c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the user ID from the request context.
		userID, found := middleware.GetUserID(r)
		if !found {
			// return a no content status if the user ID is not found.
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// Get the list of URLs created by user.
		urls, err := repo.GetByUser(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(urls) == 0 {
			// Return a no content status if the user has no URLs.
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// Preallocate a slice of URLItem objects from the URLs.
		URLItems := make([]models.URLItem, len(urls))
		for i, v := range urls {
			// Build the full shortened url from new id and service URL.
			shortURL := strings.Join([]string{c.BaseURL, v.ShortURL}, "/")
			URLItems[i] = models.URLItem{
				ShortURL:    shortURL,
				OriginalURL: v.LongURL,
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(URLItems); err != nil {
			http.Error(w, "error building the response", http.StatusInternalServerError)
			return
		}
	}
}
