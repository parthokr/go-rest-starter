package middleware

import (
	"log"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// do auth
		log.Printf("Authenticating user")
		jwt := r.Header.Get("Authorization")
		if jwt == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
