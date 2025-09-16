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

	conn := db.ConnectDB()
	mux := http.NewServeMux()

	mux.HandleFunc("/health", health.Handler)

	userService := user.NewService(conn)
	walletService := wallet.NewService(conn)

	authHandler := auth.NewHandler(userService, walletService)
	walletHandler := wallet.NewHandler(walletService)

	mux.HandleFunc("/auth/signup", authHandler.Signup)
	mux.HandleFunc("/auth/login", authHandler.Login)

	mux.Handle("/wallet/balance", auth.JWTMiddleware(http.HandlerFunc(walletHandler.GetBalance)))
	mux.Handle("/wallet/topup", auth.JWTMiddleware(http.HandlerFunc(walletHandler.TopUp)))
	mux.Handle("/wallet/transactions", auth.JWTMiddleware(http.HandlerFunc(walletHandler.GetTransactions)))
	mux.Handle("/wallet/withdraw", auth.JWTMiddleware(http.HandlerFunc(walletHandler.Withdraw)))

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Devbox API running on %s...\n", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func main() {
	startServer()
}
