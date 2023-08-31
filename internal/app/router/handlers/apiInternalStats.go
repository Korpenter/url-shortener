package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Mldlr/url-shortener/internal/app/service"
)

// APIInternalStats returns the amount of registered users and stored urls
func APIInternalStats(shortener service.ShortenerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats, err := shortener.Stats(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			http.Error(w, "error building the response", http.StatusInternalServerError)
			return
		}
	}
}
