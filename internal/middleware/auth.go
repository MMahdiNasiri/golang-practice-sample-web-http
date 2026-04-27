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
				handleUnauthorized(w, r)
				return
			}
			token := strings.TrimPrefix(tokenHeader, "Bearer ")

			userId, err := authService.ValidateToken(r.Context(), token)
			if err != nil {
				handleUnauthorized(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userId)
			next(w, r.WithContext(ctx))
		}
	}
}

func handleUnauthorized(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, "/signin-page", http.StatusFound)
}
