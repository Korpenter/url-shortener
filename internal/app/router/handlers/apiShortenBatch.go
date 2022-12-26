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

func APIShortenBatch(repo storage.Repository, c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bodyItems []models.BatchReqItem
		if err := json.NewDecoder(r.Body).Decode(&bodyItems); err != nil {
			http.Error(w, "error reading request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		urls := make(map[string]*models.URL, len(bodyItems))
		respItems := make([]models.BatchRespItem, len(bodyItems))
		for i, v := range bodyItems {
			if !validators.IsURL(v.OrigURL) {
				respItems[i] = models.BatchRespItem{
					CorID:    v.CorID,
					ShortURL: "incorrect url",
				}
				continue
			}
			id, err := repo.NewID(v.OrigURL)
			if err != nil {
				http.Error(w, fmt.Sprintf("error getting new id: %v", err), http.StatusInternalServerError)
				return
			}
			userID, found := middleware.GetUserID(r)
			if !found {
				http.Error(w, fmt.Sprintf("error getting user cookie: %v", err), http.StatusInternalServerError)
				return
			}
			urls[v.CorID] = &models.URL{
				ShortURL: id,
				LongURL:  v.OrigURL,
				UserID:   userID,
			}
			respItems[i] = models.BatchRespItem{
				CorID:    v.CorID,
				ShortURL: strings.Join([]string{c.BaseURL, id}, "/"),
			}
		}
		duplicates, err := repo.AddBatch(r.Context(), urls)
		if err != nil {
			http.Error(w, fmt.Sprintf("error adding record to db: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if duplicates {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		if err := json.NewEncoder(w).Encode(respItems); err != nil || len(respItems) == 0 {
			http.Error(w, "error building the response", http.StatusInternalServerError)
			return
		}
	}
}
