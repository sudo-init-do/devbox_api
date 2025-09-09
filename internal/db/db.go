package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func ConnectDB() *sql.DB {
	// Always load .env for local/dev
	_ = godotenv.Load()

	appEnv := os.Getenv("APP_ENV")

	var dbURL string
	if appEnv == "docker" {
		dbURL = os.Getenv("DATABASE_URL_DOCKER")
	} else {
		dbURL = os.Getenv("DATABASE_URL_LOCAL")
	}

	if dbURL == "" {
		panic("❌ DATABASE_URL is not set for env: " + appEnv)
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(fmt.Sprintf("❌ Failed to connect to DB: %v", err))
	}

	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("❌ Failed to ping DB: %v", err))
	}

	log.Println("✅ Connected to Postgres at", dbURL)
	return db
}
