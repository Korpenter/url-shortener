package handlers

import (
	"compress/gzip"
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"github.com/Mldlr/url-shortener/internal/app/utils/validators"
	"io"
	"log"
	"net/http"
)

// Shorten returns a handler that shortens links and adds them to db
func Shorten(repo storage.Repository, c *config.Config) http.HandlerFunc {
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
		b, err := io.ReadAll(reader)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request %v :", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		long := string(b)
		if !validators.IsURL(long) {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}
		id, err := repo.NewID()
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting new id: %v", err), http.StatusInternalServerError)
			return
		}
		id62 := encoders.ToRBase62(id)
		short, err := repo.Add(long, id62)
		if err != nil {
			http.Error(w, fmt.Sprintf("error adding record to db: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		if _, err = io.WriteString(w, c.BaseURL+"/"+short); err != nil {
			log.Println(err)
		}
	}
}
