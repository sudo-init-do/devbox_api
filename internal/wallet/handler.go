package wallet

import (
	"encoding/json"
	"net/http"

	"github.com/sudo-init-do/devbox_api/internal/auth"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized: missing user ID", http.StatusUnauthorized)
		return
	}

	balance, err := h.service.GetBalance(userID)
	if err != nil {
		http.Error(w, "Failed to fetch balance", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userID,
		"balance": balance,
	})
}
