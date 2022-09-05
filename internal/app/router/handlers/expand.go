package handlers

import (
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"net/http"
	"strings"
)

func Expand(repo storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.Split(r.URL.Path, "/")[1:]
		l, err := repo.Get(id[0])
		if err != nil {
			http.Error(w, "invalid id", http.StatusNotFound)
			return
		}
		w.Header().Set("Location", l)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
