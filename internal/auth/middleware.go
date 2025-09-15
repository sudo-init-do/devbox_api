package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	userIDKey contextKey = "userID"
	roleKey   contextKey = "role"
)

// JWTMiddleware validates JWT tokens and attaches user data to the request context
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			fmt.Println("Missing Authorization header")
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Expect format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			fmt.Println("Invalid Authorization format:", authHeader)
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		fmt.Println("🔑 Raw token received:", tokenString)

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			fmt.Println("⚠️ JWT_SECRET not set in environment")
		}

		// Parse and validate JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Println(" Unexpected signing method:", token.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil {
			fmt.Println("JWT parse/validation error:", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println("JWT claims decoded:", claims)

			userID, _ := claims["user_id"].(string)
			role, _ := claims["role"].(string)

			// Add values to request context
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			ctx = context.WithValue(ctx, roleKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		fmt.Println("Invalid or missing claims in token")
		http.Error(w, "Invalid token", http.StatusUnauthorized)
	})
}

// GetUserID extracts the user ID from request context
func GetUserID(r *http.Request) string {
	if v := r.Context().Value(userIDKey); v != nil {
		if userID, ok := v.(string); ok {
			return userID
		}
	}
	return ""
}

// GetUserRole extracts the role from request context
func GetUserRole(r *http.Request) string {
	if v := r.Context().Value(roleKey); v != nil {
		if role, ok := v.(string); ok {
			return role
		}
	}
	return ""
}
