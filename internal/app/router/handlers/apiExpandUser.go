package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/router/middleware"
	"net/http"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/model"
	"github.com/Mldlr/url-shortener/internal/app/storage"
)

func APIUserExpand(repo storage.Repository, c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, found := middleware.GetUserID(r)
		switch found {
		case true:
			urls, err := repo.GetByUser(userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if len(urls) == 0 {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			var URLItems []model.URLItem
			for _, v := range urls {
				shortURL := fmt.Sprintf("%s/%s", c.BaseURL, v.ShortURL)
				URLItems = append(URLItems, model.URLItem{
					ShortURL:    shortURL,
					OriginalURL: v.LongURL,
				})
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(URLItems); err != nil {
				http.Error(w, "error building the response", http.StatusInternalServerError)
				return
			}
		case false:
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
}
