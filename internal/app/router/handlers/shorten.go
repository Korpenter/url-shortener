package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"github.com/Mldlr/url-shortener/internal/app/utils/validators"
)

// Shorten returns a handler that shortens links and adds them to db
func Shorten(repo storage.Repository, c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error reading request", http.StatusBadRequest)
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
		userID, _ := r.Cookie("user_id")
		short, err := repo.Add(long, id62, userID.Value)
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
