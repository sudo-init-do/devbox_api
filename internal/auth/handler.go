package auth

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/sudo-init-do/devbox_api/internal/user"
	"github.com/sudo-init-do/devbox_api/internal/wallet"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	userService   *user.Service
	walletService *wallet.Service
}

func NewHandler(userService *user.Service, walletService *wallet.Service) *Handler {
	return &Handler{
		userService:   userService,
		walletService: walletService,
	}
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Signup handles new user registration + wallet creation
func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	u, err := h.userService.Create(req.Email, string(hashed), "creator")
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	if err := h.walletService.CreateWallet(u.ID.String()); err != nil {
		log.Printf("[AuthHandler] Failed to create wallet for user_id=%s: %v", u.ID.String(), err)
		http.Error(w, "Failed to create wallet", http.StatusInternalServerError)
		return
	}

	token, _ := GenerateToken(u.ID.String(), u.Role)
	json.NewEncoder(w).Encode(SignupResponse{Token: token})
}

// Login authenticates user
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	u, err := h.userService.FindByEmail(req.Email)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, _ := GenerateToken(u.ID.String(), u.Role)
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}
