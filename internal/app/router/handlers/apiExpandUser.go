package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/service"
	"github.com/Mldlr/url-shortener/internal/app/utils/helpers"
)

// APIUserExpand retrieves the list of shortened URLs
// created by a user and returns them as a JSON array.
func APIUserExpand(shortener service.ShortenerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the user ID from the request context.
		userID, found := helpers.GetUserID(r)
		if !found {
			// return a no content status if the user ID is not found.
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// Get the list of URLs created by user.
		urls, err := shortener.ExpandUser(r.Context(), userID)
		if err != nil {
			if !errors.Is(err, models.ErrNoContent) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Error(w, err.Error(), http.StatusNoContent)
			return
		}
		URLItems := make([]*models.URLItem, len(urls))
		for i, v := range urls {
			// Build the full shortened url from new id and service URL.
			URLItems[i] = &models.URLItem{
				ShortURL:    shortener.BuildURL(v.ShortURL),
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
