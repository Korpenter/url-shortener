package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/router/middleware"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/Mldlr/url-shortener/internal/app/utils/validators"
	"net/http"
	"strings"
)

// APIShortenBatch shortens a batch of URLs.
func APIShortenBatch(repo storage.Repository, c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode the request body as a slice of BatchReqItem models.
		var bodyItems []models.BatchReqItem
		if err := json.NewDecoder(r.Body).Decode(&bodyItems); err != nil {
			http.Error(w, "error reading request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Preallocate maps to store the URLs and response items.
		urls := make(map[string]*models.URL, len(bodyItems))
		respItems := make([]models.BatchRespItem, len(bodyItems))
		for i, v := range bodyItems {
			// Check if the original URL is valid.
			if !validators.IsURL(v.OrigURL) {
				// If the URL is not valid, set the response item to indicate a bad URL request.
				respItems[i] = models.BatchRespItem{
					CorID:    v.CorID,
					ShortURL: "incorrect url",
				}
				continue
			}
			// Generate a short ID for the URL.
			id, err := repo.NewID(v.OrigURL)
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
			// Create a new URL model and add it to the URLs map.
			urls[v.CorID] = &models.URL{
				ShortURL: id,
				LongURL:  v.OrigURL,
				UserID:   userID,
			}
			// Create a new response item.
			respItems[i] = models.BatchRespItem{
				CorID:    v.CorID,
				ShortURL: strings.Join([]string{c.BaseURL, id}, "/"),
			}
		}
		// Add the URLs to the repository.
		duplicates, err := repo.AddBatch(r.Context(), urls)
		if err != nil {
			http.Error(w, fmt.Sprintf("error adding record to db: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if duplicates {
			// If any of the URLs already exist, return a conflict status.
			w.WriteHeader(http.StatusConflict)
		} else {
			// If all URLs are new, return a created status.
			w.WriteHeader(http.StatusCreated)
		}
		if err := json.NewEncoder(w).Encode(respItems); err != nil || len(respItems) == 0 {
			http.Error(w, "error building the response", http.StatusInternalServerError)
			return
		}
	}
}
