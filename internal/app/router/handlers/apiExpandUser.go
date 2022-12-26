package handlers

import (
	"encoding/json"
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/router/middleware"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"net/http"
	"strings"
)

func APIUserExpand(repo storage.Repository, c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, found := middleware.GetUserID(r)
		if !found {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		urls, err := repo.GetByUser(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(urls) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		URLItems := make([]models.URLItem, len(urls))
		for i, v := range urls {
			shortURL := strings.Join([]string{c.BaseURL, v.ShortURL}, "")
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
