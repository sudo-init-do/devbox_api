package wallet

import (
	"database/sql"
	"fmt"
)

type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) CreateWallet(userID string) error {
	query := `
		INSERT INTO wallets (user_id, balance)
		VALUES ($1, 0)
	`
	_, err := s.DB.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}
	return nil
}
