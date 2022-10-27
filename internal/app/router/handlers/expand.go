package handlers

import (
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"net/http"
	"strings"
)

// Expand returns a handler that gets original link from db
func Expand(repo storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.Split(r.URL.Path, "/")[1:]
		l, err := repo.Get(id[0], r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Location", l)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
