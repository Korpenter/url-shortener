package handlers

import (
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
		b, err := io.ReadAll(r.Body)
		defer func() {
			if err = r.Body.Close(); err != nil {
				log.Println(err)
			}
		}()
		if err != nil {
			http.Error(w, "error reading request", http.StatusBadRequest)
			return
		}
		long := string(b)
		if !validators.IsURL(long) {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}
		id := encoders.ToRBase62(repo.NewID())
		short := repo.Add(long, id)
		w.Header().Set("Content-Type", "text/plain;")
		w.WriteHeader(http.StatusCreated)
		if _, err = io.WriteString(w, c.BaseURL+"/"+short); err != nil {
			log.Println(err)
		}
	}
}
