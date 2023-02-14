package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/service"
	"github.com/Mldlr/url-shortener/internal/app/utils/helpers"
)

// APIShortenBatch shortens a batch of URLs.
func APIShortenBatch(shortener service.ShortenerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the user ID from the request.
		userID, found := helpers.GetUserID(r)
		if !found {
			http.Error(w, "error getting user cookie", http.StatusInternalServerError)
			return
		}
		// Decode the request body as a slice of BatchReqItem models.
		var bodyItems []models.BatchReqItem
		if err := json.NewDecoder(r.Body).Decode(&bodyItems); err != nil {
			http.Error(w, "error reading request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Preallocate maps to store the URLs and response items.
		urls := make([]*models.URL, len(bodyItems))
		for i, v := range bodyItems {
			// Create a new URL model and add it to the URLs map.
			urls[i] = &models.URL{
				LongURL: v.OrigURL,
				UserID:  userID,
			}
		}
		// Add the URLs to the repository.
		var statusCode int
		shortenedURLs, err := shortener.ShortenBatch(r.Context(), userID, urls)
		if err != nil {
			// If there is an error, and its not a duplicate url
			if !errors.Is(err, models.ErrDuplicate) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// If it is a duplicate url
			statusCode = http.StatusConflict
		} else {
			// If URL is new, return a created status.
			statusCode = http.StatusCreated
		}
		respItems := make([]models.BatchRespItem, len(bodyItems))
		for i, v := range bodyItems {
			// Create a new response item.
			respItems[i] = models.BatchRespItem{
				CorID:    v.CorID,
				ShortURL: shortener.BuildURL(shortenedURLs[i].ShortURL),
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if err := json.NewEncoder(w).Encode(respItems); err != nil || len(respItems) == 0 {
			http.Error(w, "error building the response", http.StatusInternalServerError)
			return
		}
	}
}
