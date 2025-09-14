package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sudo-init-do/devbox_api/internal/auth"
	"github.com/sudo-init-do/devbox_api/internal/db"
	"github.com/sudo-init-do/devbox_api/internal/health"
	"github.com/sudo-init-do/devbox_api/internal/user"
	"github.com/sudo-init-do/devbox_api/internal/wallet"
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

	walletService := wallet.NewService(conn)
	walletHandler := wallet.NewHandler(walletService)

	// Auth routes
	mux.HandleFunc("/auth/signup", authHandler.Signup)
	mux.HandleFunc("/auth/login", authHandler.Login)

	// Wallet routes (protected by JWT)
	mux.Handle("/wallet/balance", auth.JWTMiddleware(http.HandlerFunc(walletHandler.GetBalance)))

	addr := fmt.Sprintf(":%s", port)
	log.Printf(" Devbox API running on %s...\n", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func main() {
	startServer()
}
