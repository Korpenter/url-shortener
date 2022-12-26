package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/router/middleware"
	"net/http"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/models"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/Mldlr/url-shortener/internal/app/utils/validators"
)

// Expand returns a handler that gets original link from db
func APIShorten(repo storage.Repository, c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body *models.URL
		var err error
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "error reading request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		if !validators.IsURL(body.LongURL) {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}
		body.ShortURL, err = repo.NewID(body.LongURL)
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting new id: %v", err), http.StatusInternalServerError)
			return
		}
		userID, found := middleware.GetUserID(r)
		if !found {
			http.Error(w, fmt.Sprintf("error getting user cookie: %v", err), http.StatusInternalServerError)
			return
		}
		body.UserID = userID
		duplicates, err := repo.Add(r.Context(), body)
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
		if err := json.NewEncoder(w).Encode(models.Response{Result: c.BaseURL + body.ShortURL}); err != nil {
			http.Error(w, "error building the response", http.StatusInternalServerError)
			return
		}
	}
}
