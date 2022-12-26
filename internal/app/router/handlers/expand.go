package handlers

import (
	"net/http"
	"strings"

	"github.com/Mldlr/url-shortener/internal/app/storage"
)

// Expand returns a handler that gets original link from db
func Expand(repo storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.Split(r.URL.Path, "/")[1:]
		url, err := repo.Get(r.Context(), id[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if url.Deleted {
			w.WriteHeader(http.StatusGone)
			return
		}
		w.Header().Set("Location", url.LongURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
