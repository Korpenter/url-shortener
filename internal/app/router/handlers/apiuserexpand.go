package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/model"
	"github.com/Mldlr/url-shortener/internal/app/storage"
)

func APIUserExpand(repo storage.Repository, c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Cookie("user_id")
		urls, _ := repo.GetByUser(user.Value)
		fmt.Println(urls)
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
	}
}
