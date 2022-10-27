package handlers

import (
	"github.com/Mldlr/url-shortener/internal/app/storage"
	"net/http"
)

func Ping(repo storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := repo.Ping(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
