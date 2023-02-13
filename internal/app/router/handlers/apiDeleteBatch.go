// Package handlers provides HTTP handlers for the url-shortener.
package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Mldlr/url-shortener/internal/app/service"
	"github.com/Mldlr/url-shortener/internal/app/utils/helpers"
)

// APIDeleteBatch processes a batch request to delete multiple shortened URLs.
func APIDeleteBatch(shortener service.ShortenerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the user ID from the request context.
		userID, found := helpers.GetUserID(r)
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
		shortener.APIDeleteBatch(urlIDs, userID)
		// Return an accepted status to indicate that the request has been
		// received and is being processed.
		w.WriteHeader(http.StatusAccepted)
	}
}
