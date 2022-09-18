package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"github.com/Mldlr/url-shortener/internal/app/utils/validators"
	"io"
	"net/http"
)

type Request struct {
	URL string `json:"url,omitempty"`
}

type Response struct {
	Result string `json:"result"`
}

// Expand returns a handler that gets original link from db
func APIShorten(repo storage.Repository, c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reader io.ReadCloser
		var err error
		switch r.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(r.Body)
			defer reader.Close()
		default:
			reader = r.Body
		}
		var body *Request
		if err := json.NewDecoder(reader).Decode(&body); err != nil {
			http.Error(w, "error reading request", http.StatusBadRequest)
		}
		defer r.Body.Close()
		if !validators.IsURL(body.URL) {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}
		id, err := repo.NewID()
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting new id: %v", err), http.StatusInternalServerError)
			return
		}
		id62 := encoders.ToRBase62(id)
		short, err := repo.Add(body.URL, id62)
		if err != nil {
			http.Error(w, fmt.Sprintf("error adding record to db: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(Response{c.BaseURL + "/" + short}); err != nil {
			http.Error(w, "error building the response", http.StatusInternalServerError)
			return
		}
	}
}
