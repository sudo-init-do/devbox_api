package wallet

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sudo-init-do/devbox_api/internal/auth"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// GET /wallet/balance
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	log.Printf("[WalletHandler] /wallet/balance called by user_id=%s", userID)

	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	balance, err := h.service.GetBalance(userID)
	if err != nil {
		log.Printf("[WalletHandler] Failed to fetch balance for user_id=%s: %v", userID, err)
		http.Error(w, "Failed to fetch balance", http.StatusInternalServerError)
		return
	}

	log.Printf("[WalletHandler] Balance fetched successfully for user_id=%s: balance=%d", userID, balance)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userID,
		"balance": balance,
	})
}

// POST /wallet/topup
func (h *Handler) TopUp(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	log.Printf("[WalletHandler] /wallet/topup called by user_id=%s", userID)

	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type TopUpRequest struct {
		Amount    int    `json:"amount"`
		Reference string `json:"reference"`
	}
	var req TopUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[WalletHandler] Invalid request body for user_id=%s: %v", userID, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 || req.Reference == "" {
		log.Printf("[WalletHandler] Invalid top-up request for user_id=%s: amount=%d reference=%s",
			userID, req.Amount, req.Reference)
		http.Error(w, "Invalid top-up request", http.StatusBadRequest)
		return
	}

	tx, balance, err := h.service.TopUp(userID, int64(req.Amount), req.Reference)
	if err != nil {
		log.Printf("[WalletHandler] Top-up failed for user_id=%s: %v", userID, err)
		http.Error(w, "Failed to process top-up", http.StatusInternalServerError)
		return
	}

	log.Printf("[WalletHandler] Top-up success for user_id=%s: new_balance=%d, tx_id=%s",
		userID, balance, tx.ID)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":     userID,
		"transaction": tx,
		"balance":     balance,
	})
}

// GET /wallet/transactions
func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	log.Printf("[WalletHandler] /wallet/transactions called by user_id=%s", userID)

	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	txs, err := h.service.GetTransactions(userID)
	if err != nil {
		log.Printf("[WalletHandler] Failed to fetch transactions for user_id=%s: %v", userID, err)
		http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
		return
	}

	log.Printf("[WalletHandler] Transactions fetched successfully for user_id=%s, count=%d", userID, len(txs))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":      userID,
		"transactions": txs,
	})
}
