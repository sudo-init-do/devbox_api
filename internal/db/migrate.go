package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// RunMigrations runs migrations from the migrations.sql file
func RunMigrations() {
	conn := ConnectDB()
	defer conn.Close()

	// Find migrations.sql (root/migrate/migrations.sql)
	path := filepath.Join("migrate", "migrations.sql")
	sqlBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("❌ Failed to read migrations file: %v", err)
	}

	_, err = conn.Exec(string(sqlBytes))
	if err != nil {
		log.Fatalf("❌ Failed to execute migrations: %v", err)
	}

	fmt.Println("✅ Migrations applied successfully")
}
