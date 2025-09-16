package wallet

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// CreateWallet creates a wallet for a new user if not exists
func (s *Service) CreateWallet(userID string) error {
	_, err := s.db.Exec(`
		INSERT INTO wallets (id, user_id, balance, created_at, updated_at)
		VALUES ($1, $2, 0, NOW(), NOW())
		ON CONFLICT (user_id) DO NOTHING
	`, uuid.New().String(), userID)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}
	return nil
}

// GetBalance returns the user’s wallet balance
func (s *Service) GetBalance(userID string) (int64, error) {
	var balance int64
	err := s.db.QueryRow("SELECT balance FROM wallets WHERE user_id=$1", userID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch balance: %w", err)
	}
	return balance, nil
}

// TopUp adds coins to user wallet
func (s *Service) TopUp(userID string, amount int64, reference string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE wallets SET balance = balance + $1, updated_at = NOW()
		WHERE user_id = $2
	`, amount, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO transactions (id, user_id, amount, type, reference, created_at)
		VALUES ($1, $2, $3, 'topup', $4, $5)
	`, uuid.New().String(), userID, amount, reference, time.Now())
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Withdraw deducts coins from wallet
func (s *Service) Withdraw(userID string, amount int64, reference string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// check balance
	var balance int64
	err = tx.QueryRow("SELECT balance FROM wallets WHERE user_id=$1", userID).Scan(&balance)
	if err != nil {
		return err
	}
	if balance < amount {
		return fmt.Errorf("insufficient balance")
	}

	_, err = tx.Exec(`
		UPDATE wallets SET balance = balance - $1, updated_at = NOW()
		WHERE user_id = $2
	`, amount, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO transactions (id, user_id, amount, type, reference, created_at)
		VALUES ($1, $2, $3, 'withdraw', $4, $5)
	`, uuid.New().String(), userID, amount, reference, time.Now())
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetTransactions returns all wallet transactions
func (s *Service) GetTransactions(userID string) ([]map[string]interface{}, error) {
	rows, err := s.db.Query(`
		SELECT id, amount, type, reference, created_at
		FROM transactions WHERE user_id=$1 ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []map[string]interface{}
	for rows.Next() {
		var id, ttype, ref string
		var amount int64
		var created time.Time

		if err := rows.Scan(&id, &amount, &ttype, &ref, &created); err != nil {
			return nil, err
		}

		txs = append(txs, map[string]interface{}{
			"id":        id,
			"amount":    amount,
			"type":      ttype,
			"reference": ref,
			"created":   created,
		})
	}
	return txs, nil
}
