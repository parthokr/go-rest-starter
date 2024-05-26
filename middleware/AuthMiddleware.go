package middleware

import (
	"context"
	"go-rest-starter/server"
	"go-rest-starter/utils"
	"log"
	"net/http"
	"strings"
	"time"
)

func AuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// do auth
		log.Printf("Authenticating user")
		token := r.Header.Get("Authorization")
		if token == "" {
			errResp := server.NewAPIError(
				"NO_AUTHORIZATION_HEADER",
				r.URL.Path,
				map[string]string{"token": "Bearer token is required"},
				time.Now().Unix())
			err := server.WriteJSON(w, http.StatusUnauthorized, errResp)
			if err != nil {
				log.Println("Error writing JSON response")
			}
			return
		}
		bearerToken := strings.Split(token, "Bearer ")
		if len(bearerToken) != 2 {
			errResp := server.NewAPIError(
				"INVALID_TOKEN_FORMAT",
				r.URL.Path,
				map[string]string{"token": "Invalid token format"},
				time.Now().Unix())
			err := server.WriteJSON(w, http.StatusUnauthorized, errResp)
			if err != nil {
				log.Println("Error writing JSON response")
			}
			return
		}
		// check if token is valid
		parsedToken, err := utils.VerifyToken(bearerToken[1])
		if err != nil {
			err = server.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			if err != nil {
				log.Println("Error writing JSON response")
			}
			return
		}
		log.Printf("Token: %+v", parsedToken)
		// check if token is expired
		iat := parsedToken["iat"].(float64)
		isExpired := time.Since(time.Unix(int64(iat), 0)) > 30*time.Second
		if isExpired {
			errResp := server.NewAPIError(
				"TOKEN_EXPIRED",
				r.URL.Path,
				map[string]string{"token": "Token has been expired"},
				time.Now().Unix())
			err = server.WriteJSON(w, http.StatusUnauthorized, errResp)
			if err != nil {
				log.Println("Error writing JSON response")
			}
			return
		}
		// update request context with user info
		r = r.Clone(context.WithValue(r.Context(), "username", parsedToken["username"]))
		next.ServeHTTP(w, r)
	}
}
