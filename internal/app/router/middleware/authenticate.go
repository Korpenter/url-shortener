package middleware

import (
	"net/http"
	"time"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"github.com/google/uuid"
)

type Auth struct {
	Config *config.Config
}

func (a Auth) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := r.Cookie("user_id")
		switch err {
		case nil:
			userID := user.Value
			requestSignature, _ := r.Cookie("signature")
			if requestSignature.Value == encoders.HMACString(userID, a.Config.SecretKey) {
				break
			}
			fallthrough
		default:
			user = &http.Cookie{
				Name:    "user_id",
				Path:    "/",
				Expires: time.Now().Add(time.Hour * 24 * 7),
				Value:   uuid.New().String(),
			}
			signature := &http.Cookie{
				Name:    "signature",
				Path:    "/",
				Expires: time.Now().Add(time.Hour * 24 * 7),
				Value:   encoders.HMACString(uuid.New().String(), a.Config.SecretKey),
			}
			r.AddCookie(user)
			r.AddCookie(signature)
			http.SetCookie(w, user)
			http.SetCookie(w, signature)
		}
		next.ServeHTTP(w, r)
	})
}
