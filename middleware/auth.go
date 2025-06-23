package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// BasicAuth middleware for protecting admin endpoints
func BasicAuth(username, password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				w.Header().Set("WWW-Authenticate", `Basic realm="Admin Area"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(auth, "Basic ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			payload, err := base64.StdEncoding.DecodeString(auth[6:])
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			pair := strings.SplitN(string(payload), ":", 2)
			if len(pair) != 2 || pair[0] != username || pair[1] != password {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}