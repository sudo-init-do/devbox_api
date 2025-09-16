package wallet

import (
	"encoding/json"
	"net/http"

	"github.com/sudo-init-do/devbox_api/internal/auth"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GET /wallet/balance
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r) // ✅ call from auth middleware
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	balance, err := h.service.GetBalance(userID)
	if err != nil {
		http.Error(w, "Failed to fetch balance", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"balance": balance,
	})
}

// POST /wallet/topup
func (h *Handler) TopUp(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Amount    int64  `json:"amount"`
		Reference string `json:"reference"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.TopUp(userID, req.Amount, req.Reference); err != nil {
		http.Error(w, "Failed to top up wallet", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Top-up successful",
	})
}

// POST /wallet/withdraw
func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Amount    int64  `json:"amount"`
		Reference string `json:"reference"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.Withdraw(userID, req.Amount, req.Reference); err != nil {
		http.Error(w, "Failed to withdraw", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Withdrawal successful",
	})
}

// GET /wallet/transactions
func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	txs, err := h.service.GetTransactions(userID)
	if err != nil {
		http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(txs)
}
