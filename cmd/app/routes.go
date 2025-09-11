package main

import (
	"database/sql"
	"net/http"

	"github.com/sudo-init-do/devbox_api/internal/auth"
	"github.com/sudo-init-do/devbox_api/internal/health"
)

func initRoutes(mux *http.ServeMux, conn *sql.DB) {
	// Health check route
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		// redirect to /health
		http.Redirect(writer, request, "/health", http.StatusMovedPermanently)
	})
	mux.HandleFunc("/health", health.Handler)

	// Auth routes
	// TODO: Passing db connection directly to handlers is not ideal as it creates tight coupling
	// and requires the connection to be instantiated/available every time a handler is called.
	// Consider using dependency injection, a service layer, or connection pooling pattern instead.
	mux.Handle("/auth/signup", auth.SignupHandler(conn))
	mux.Handle("/auth/login", auth.LoginHandler(conn))
}
