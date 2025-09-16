package wallet

import "time"

// Wallet model
type Wallet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Transaction model
type Transaction struct {
	ID        string    `json:"id"`
	WalletID  string    `json:"wallet_id"`
	Amount    int64     `json:"amount"`
	Type      string    `json:"type"`      // credit or debit
	Reference string    `json:"reference"` // unique reference
	Status    string    `json:"status"`    // completed or pending
	CreatedAt time.Time `json:"created_at"`
}
