package db

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func Seed(db *sql.DB) {
	log.Println("🌱 Running seed data...")

	// Hash password
	password := "password123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Insert user
	var userID string
	err := db.QueryRow(`
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO NOTHING
		RETURNING id
	`, "test@example.com", string(hash), "creator").Scan(&userID)

	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("❌ Failed to insert user: %v", err)
	}

	// Insert wallet if user was created
	if userID != "" {
		_, err = db.Exec(`
			INSERT INTO wallets (user_id, balance)
			VALUES ($1, $2)
		`, userID, 0)
		if err != nil {
			log.Fatalf("❌ Failed to insert wallet: %v", err)
		}
	}

	log.Println("✅ Seeding complete!")
}
