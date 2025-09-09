package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sudo-init-do/devbox_api/internal/db"
	"github.com/sudo-init-do/devbox_api/internal/health"
)

func main() {
	// Check if a CLI command was passed
	if len(os.Args) > 1 {
		cmd := os.Args[1]

		switch cmd {
		case "migrate":
			db.RunMigrations()
			return
		case "seed":
			conn := db.ConnectDB()
			db.Seed(conn)
			return
		default:
			log.Fatalf("❌ Unknown command: %s", cmd)
		}
	}

	// Default: run server
	startServer()
}

func startServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", health.Handler)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Devbox API running on %s...\n", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
