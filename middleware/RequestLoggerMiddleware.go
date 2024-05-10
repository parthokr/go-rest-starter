package middleware

import (
	"log"
	"net/http"
)

func RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method %s path: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}
