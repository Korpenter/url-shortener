// Package middleware provides custom middleware.
package middleware

import (
	"net/http"
	"time"

	"github.com/Mldlr/url-shortener/internal/app/config"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	"github.com/google/uuid"
)

// Auth is an authentication middleware.
type Auth struct {
	Config *config.Config
}

// Authenticate authenticates user request by adding user ID and signature cookies to the request
// if they are not present or invalid.
func (a Auth) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the user ID from the request context.
		switch userID, found := GetUserID(r); found {
		case true:
			// If the ID is found, check the signature.
			requestSignature, err := r.Cookie("signature")
			if err == nil && requestSignature.Value == encoders.HMACString(userID, a.Config.SecretKey) {
				break
			}
			fallthrough
		default:
			// If the ID is not found, or the signature is invalid, create new ID and signature.
			id := uuid.New().String()
			userID := &http.Cookie{
				Name:    "user_id",
				Path:    "/",
				Expires: time.Now().Add(time.Hour * 24 * 7),
				Value:   id,
			}
			signature := &http.Cookie{
				Name:    "signature",
				Path:    "/",
				Expires: time.Now().Add(time.Hour * 24 * 7),
				Value:   encoders.HMACString(id, a.Config.SecretKey),
			}
			// Sign the user request.
			r.AddCookie(userID)
			r.AddCookie(signature)
			http.SetCookie(w, userID)
			http.SetCookie(w, signature)
		}
		next.ServeHTTP(w, r)
	})
}

// GetUserID check for user_id cookie in request and returns its value if it is present.
func GetUserID(r *http.Request) (string, bool) {
	userID, err := r.Cookie("user_id")
	if err != nil {
		return "", false
	}
	return userID.Value, true
}
