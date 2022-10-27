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
		switch userID, found := GetUserID(r); found {
		case true:
			requestSignature, _ := r.Cookie("signature")
			if requestSignature.Value == encoders.HMACString(userID, a.Config.SecretKey) {
				break
			}
			fallthrough
		default:
			userID := &http.Cookie{
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
			r.AddCookie(userID)
			r.AddCookie(signature)
			http.SetCookie(w, userID)
			http.SetCookie(w, signature)
		}
		next.ServeHTTP(w, r)
	})
}

func GetUserID(r *http.Request) (string, bool) {
	userID, err := r.Cookie("user_id")
	if err != nil {
		return "", false
	}
	return userID.Value, true
}
