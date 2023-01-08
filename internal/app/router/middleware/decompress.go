package middleware

import (
	"compress/gzip"
	"net/http"
)

// Decompress decompresses the request body if it is gzipped.
func Decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Header.Get("Content-Encoding") {
		case "gzip":
			// Read request body.
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer reader.Close()
			// Get un-gzipped body.
			r.Body = reader
			next.ServeHTTP(w, r)
		default:
			next.ServeHTTP(w, r)
		}
	})
}
