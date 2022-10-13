package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/model"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"github.com/Mldlr/url-shortener/internal/app/utils/validators"
	"net/http"
)

func APIShortenBatch(repo storage.Repository, c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bodyItems []model.BatchReqItem
		var respItems []model.BatchRespItem
		var urls []*model.URL
		var url model.URL
		if err := json.NewDecoder(r.Body).Decode(&bodyItems); err != nil {
			http.Error(w, "error reading request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		for i, v := range bodyItems {
			if !validators.IsURL(v.OrigURL) {
				respItems = append(respItems, model.BatchRespItem{
					CorID:    v.CorID,
					ShortURL: "incorrect url",
				})
				continue
			}
			id, err := repo.NewID()
			if err != nil {
				http.Error(w, fmt.Sprintf("error getting new id: %v", err), http.StatusInternalServerError)
				return
			}
			userID, _ := r.Cookie("user_id")
			url = model.URL{
				ShortURL: encoders.ToRBase62(id + i),
				LongURL:  v.OrigURL,
				UserID:   userID.Value,
			}
			urls = append(urls, &url)
			respItems = append(respItems, model.BatchRespItem{
				CorID:    v.CorID,
				ShortURL: c.BaseURL + "/" + url.ShortURL,
			})
		}
		err := repo.AddBatch(urls)
		if err != nil {
			http.Error(w, fmt.Sprintf("error adding record to db: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(respItems); err != nil || len(respItems) == 0 {
			http.Error(w, "error building the response", http.StatusInternalServerError)
			return
		}
	}
}
