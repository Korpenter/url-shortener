// Package handlers provides HTTP handlers for the url-shortener.
package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/router/loader"
	"github.com/Mldlr/url-shortener/internal/app/router/middleware"
)

// APIDeleteBatch processes a batch request to delete multiple shortened URLs.
func APIDeleteBatch(loader *loader.UserLoader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the user ID from the request context.
		userID, found := middleware.GetUserID(r)
		if !found {
			http.Error(w, "error getting user cookie", http.StatusInternalServerError)
			return
		}
		// Read the request body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Unmarshal the request body into a slice of URL IDs.
		var urlIDs []string
		err = json.Unmarshal(body, &urlIDs)
		if err != nil || len(urlIDs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Create a slice of DeleteURLItem objects from the URL IDs.
		deleteURLs := make([]*models.DeleteURLItem, len(urlIDs))
		for i, v := range urlIDs {
			deleteURLs[i] = &models.DeleteURLItem{UserID: userID, ShortURL: v}
		}
		// Start a goroutine to delete the URLs asynchronously.
		go func() {
			num, err := loader.LoadAll(deleteURLs)
			if err[0] != nil {
				log.Printf("error deleing urls :%v", err[0])
			}
			var result int
			for _, v := range num {
				result += v
			}
			log.Printf("deleted %v urls", result)
		}()
		// Return an accepted status to indicate that the request has been
		// received and is being processed.
		w.WriteHeader(http.StatusAccepted)
	}
}
