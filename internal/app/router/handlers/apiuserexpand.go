package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mldlr/url-shortener/internal/app/storage"
)

func APIUserExpand(repo storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Cookie("user_id")
		urls, _ := repo.GetByUser(user.Value)
		fmt.Println(urls)
		if len(urls) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(urls); err != nil {
			http.Error(w, "error building the response", http.StatusInternalServerError)
			return
		}
	}
}
