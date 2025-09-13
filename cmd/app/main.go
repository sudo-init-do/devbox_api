package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sudo-init-do/devbox_api/internal/auth"
	"github.com/sudo-init-do/devbox_api/internal/db"
	"github.com/sudo-init-do/devbox_api/internal/health"
	"github.com/sudo-init-do/devbox_api/internal/user"
)

func startServer() {
	port := "8080"

	// Connect DB once
	conn := db.ConnectDB()

	// Setup routes
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/health", health.Handler)

	// Initialize services
	userService := user.NewService(conn)
	authHandler := auth.NewHandler(userService)

	// Auth routes
	mux.HandleFunc("/auth/signup", authHandler.Signup)
	mux.HandleFunc("/auth/login", authHandler.Login)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Devbox API running on %s...\n", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
