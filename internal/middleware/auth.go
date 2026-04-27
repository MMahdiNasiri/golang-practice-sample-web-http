package middleware

import (
	"context"
	"net/http"
	"sample-web-http/internal/authenticate"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(authService *authenticate.TokenService) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenHeader := r.Header.Get("Authorization")
			if tokenHeader == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}
			token := strings.TrimPrefix(tokenHeader, "Bearer ")

			userId, err := authService.ValidateToken(r.Context(), token)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userId)
			next(w, r.WithContext(ctx))
		}
	}
}
