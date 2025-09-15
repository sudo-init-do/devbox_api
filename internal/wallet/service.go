package wallet

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

// Service provides wallet-related operations
type Service struct {
	db *sql.DB
}

// NewService initializes wallet service
func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

type Transaction struct {
	ID        string    `json:"id"`
	WalletID  string    `json:"wallet_id"`
	Amount    int64     `json:"amount"`
	Type      string    `json:"type"`
	Reference string    `json:"reference"`
	CreatedAt time.Time `json:"created_at"`
}

// GetBalance fetches the wallet balance for a user
func (s *Service) GetBalance(userID string) (int64, error) {
	log.Printf("[WalletService] Fetching balance for user_id=%s", userID)

	var balance int64
	err := s.db.QueryRow(`SELECT balance FROM wallets WHERE user_id = $1`, userID).Scan(&balance)
	if err != nil {
		log.Printf("[WalletService] Failed to fetch balance for user_id=%s: %v", userID, err)
		return 0, err
	}

	log.Printf("[WalletService] Balance for user_id=%s: %d", userID, balance)
	return balance, nil
}

// TopUp inserts a credit transaction and updates balance
func (s *Service) TopUp(userID string, amount int64, reference string) (*Transaction, int64, error) {
	log.Printf("[WalletService] Starting top-up: user_id=%s, amount=%d, reference=%s", userID, amount, reference)

	// 1. Get wallet
	var walletID string
	err := s.db.QueryRow(`SELECT id FROM wallets WHERE user_id = $1`, userID).Scan(&walletID)
	if err != nil {
		log.Printf("[WalletService] Wallet lookup failed for user_id=%s: %v", userID, err)
		return nil, 0, err
	}
	log.Printf("[WalletService] Found wallet_id=%s for user_id=%s", walletID, userID)

	// 2. Insert transaction
	txID := uuid.New().String()
	_, err = s.db.Exec(`
		INSERT INTO transactions (id, wallet_id, amount, type, reference)
		VALUES ($1, $2, $3, $4, $5)`,
		txID, walletID, amount, "credit", reference,
	)
	if err != nil {
		log.Printf("[WalletService] Failed to insert transaction for wallet_id=%s: %v", walletID, err)
		return nil, 0, err
	}
	log.Printf("[WalletService] Inserted transaction tx_id=%s for wallet_id=%s", txID, walletID)

	// 3. Update wallet balance
	_, err = s.db.Exec(`UPDATE wallets SET balance = balance + $1 WHERE id = $2`, amount, walletID)
	if err != nil {
		log.Printf("[WalletService] Balance update failed for wallet_id=%s: %v", walletID, err)
		return nil, 0, err
	}
	log.Printf("[WalletService] Balance updated successfully for wallet_id=%s", walletID)

	// 4. Fetch new balance
	var balance int64
	err = s.db.QueryRow(`SELECT balance FROM wallets WHERE id = $1`, walletID).Scan(&balance)
	if err != nil {
		log.Printf("[WalletService] Failed to fetch updated balance for wallet_id=%s: %v", walletID, err)
		return nil, 0, err
	}
	log.Printf("[WalletService] New balance for wallet_id=%s: %d", walletID, balance)

	return &Transaction{
		ID:        txID,
		WalletID:  walletID,
		Amount:    amount,
		Type:      "credit",
		Reference: reference,
		CreatedAt: time.Now(),
	}, balance, nil
}

// GetTransactions fetches all wallet transactions for a user
func (s *Service) GetTransactions(userID string) ([]Transaction, error) {
	log.Printf("[WalletService] Fetching transactions for user_id=%s", userID)

	rows, err := s.db.Query(`
		SELECT t.id, t.wallet_id, t.amount, t.type, t.reference, t.created_at
		FROM transactions t
		JOIN wallets w ON t.wallet_id = w.id
		WHERE w.user_id = $1
		ORDER BY t.created_at DESC
	`, userID)
	if err != nil {
		log.Printf("[WalletService] Failed to fetch transactions: %v", err)
		return nil, err
	}
	defer rows.Close()

	var txs []Transaction
	for rows.Next() {
		var tx Transaction
		if err := rows.Scan(&tx.ID, &tx.WalletID, &tx.Amount, &tx.Type, &tx.Reference, &tx.CreatedAt); err != nil {
			log.Printf("[WalletService] Row scan failed: %v", err)
			return nil, err
		}
		txs = append(txs, tx)
	}

	return txs, nil
}
