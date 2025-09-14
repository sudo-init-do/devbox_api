package auth

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const (
	userIDKey ctxKey = "user_id"
	roleKey   ctxKey = "role"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userID, _ := claims["user_id"].(string)
		role, _ := claims["role"].(string)

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		ctx = context.WithValue(ctx, roleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(r *http.Request) string {
	if val, ok := r.Context().Value(userIDKey).(string); ok {
		return val
	}
	return ""
}

func GetUserRole(r *http.Request) string {
	if val, ok := r.Context().Value(roleKey).(string); ok {
		return val
	}
	return ""
}
