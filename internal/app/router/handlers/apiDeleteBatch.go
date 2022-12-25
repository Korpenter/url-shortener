package handlers

import (
	"encoding/json"
	"github.com/Mldlr/url-shortener/internal/app/model"
	"github.com/Mldlr/url-shortener/internal/app/router/loader"
	"github.com/Mldlr/url-shortener/internal/app/router/middleware"
	"io"
	"log"
	"net/http"
)

func APIDeleteBatch(loader *loader.UserLoader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, found := middleware.GetUserID(r)
		if !found {
			http.Error(w, "error getting user cookie", http.StatusInternalServerError)
			return
		}
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var urlIDs []string
		err = json.Unmarshal(body, &urlIDs)
		if err != nil || len(urlIDs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		deleteURLs := make([]*model.DeleteURLItem, 0, len(urlIDs))
		for _, v := range urlIDs {
			deleteURLs = append(deleteURLs, &model.DeleteURLItem{UserID: userID, ShortURL: v})
		}
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
		w.WriteHeader(http.StatusAccepted)
	}
}
