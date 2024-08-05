package mdl

import (
	"context"
	"net/http"
	"strings"

	jwt_service "github.com/Frozelo/music-rate-service/pkg/jwt"
)

type contextKey string

const (
	ContextKeyUserId contextKey = "userId"
)

func Auth(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt_service.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ContextKeyUserId, claims.UserId)
		nextHandler.ServeHTTP(w, r.WithContext(ctx))

	})
}
