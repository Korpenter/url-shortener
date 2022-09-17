package handlers

import (
	"encoding/json"
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"github.com/Mldlr/url-shortener/internal/app/utils/validators"
	"log"
	"net/http"
)

type Request struct {
	URL string `json:"url,omitempty"`
}

type Response struct {
	Result string `json:"result"`
}

// Expand returns a handler that gets original link from db
func ApiShorten(repo storage.Repository, c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body *Request
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "error reading request", http.StatusBadRequest)
		}
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Println(err)
			}
		}()
		if !validators.IsURL(body.URL) {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}
		id := encoders.ToRBase62(repo.NewID())
		short := repo.Add(body.URL, id)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(Response{c.BaseUrl + short}); err != nil {
			http.Error(w, "error building the response", http.StatusInternalServerError)
			return
		}
	}
}
