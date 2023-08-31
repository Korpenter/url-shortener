package helpers

import (
	"net/http"
)

// GetUserID check for user_id cookie in request and returns its value if it is present.
func GetUserID(r *http.Request) (string, bool) {
	userID, err := r.Cookie("user_id")
	if err != nil {
		return "", false
	}
	return userID.Value, true
}
