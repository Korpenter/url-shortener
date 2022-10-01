package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/model"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"github.com/Mldlr/url-shortener/internal/app/utils/validators"
)

// Expand returns a handler that gets original link from db
func APIShorten(repo storage.Repository, c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body model.URL
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "error reading request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		if !validators.IsURL(body.LongURL) {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}
		id, err := repo.NewID()
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting new id: %v", err), http.StatusInternalServerError)
			return
		}
		id62 := encoders.ToRBase62(id)
		userID, _ := r.Cookie("user_id")
		short, err := repo.Add(body.LongURL, id62, userID.Value)
		if err != nil {
			http.Error(w, fmt.Sprintf("error adding record to db: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(model.Response{c.BaseURL + "/" + short}); err != nil {
			http.Error(w, "error building the response", http.StatusInternalServerError)
			return
		}
	}
}
