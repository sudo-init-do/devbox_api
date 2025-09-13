package main

import (
	"database/sql"
	"net/http"

	"github.com/sudo-init-do/devbox_api/internal/auth"
	"github.com/sudo-init-do/devbox_api/internal/health"
	"github.com/sudo-init-do/devbox_api/internal/user"
)

func initRoutes(mux *http.ServeMux, conn *sql.DB) {
	// Health
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/health", http.StatusMovedPermanently)
	})
	mux.HandleFunc("/health", health.Handler)

	// Initialize services
	userService := user.NewService(conn)
	authHandler := auth.NewHandler(userService)

	// Auth routes
	mux.HandleFunc("/auth/signup", authHandler.Signup)
}
