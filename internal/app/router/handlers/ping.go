// Package handlers provides HTTP handlers for the url-shortener.
package handlers

import (
	"net/http"

	"github.com/Mldlr/url-shortener/internal/app/storage"
)

// Ping checks the status of the repository by using its' Ping method.
func Ping(repo storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ping the repository, cancel the operation if request is cancelled.
		err := repo.Ping(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
