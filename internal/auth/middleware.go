package auth

import (
    "net/http"
)

// Dummy middleware placeholder
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // For now, just pass through
        next.ServeHTTP(w, r)
    })
}
