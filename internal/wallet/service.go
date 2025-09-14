package wallet

import (
	"database/sql"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetBalance(userID string) (int64, error) {
	var balance int64
	err := s.db.QueryRow(`SELECT balance FROM wallets WHERE user_id = $1`, userID).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return balance, err
}
