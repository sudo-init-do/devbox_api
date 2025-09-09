package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/sudo-init-do/devbox_api/internal/user"
	"github.com/sudo-init-do/devbox_api/internal/wallet"
)

// SignupRequest holds request payload for signup
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignupResponse returns a JWT token after signup
type SignupResponse struct {
	Token string `json:"token"`
}

// LoginRequest holds request payload for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse returns a JWT token after login
type LoginResponse struct {
	Token string `json:"token"`
}

// SignupHandler handles POST /auth/signup
func SignupHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
			return
		}

		userRepo := user.NewRepository(db)
		walletService := wallet.NewService(db)

		// check if user exists
		existing, err := userRepo.FindByEmail(req.Email)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if existing != nil {
			http.Error(w, "Email already registered", http.StatusConflict)
			return
		}

		// hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// create user
		u, err := userRepo.Create(req.Email, string(hash), "creator")
		if err != nil {
			http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// create wallet
		if err := walletService.CreateWallet(u.ID); err != nil {
			http.Error(w, "Failed to create wallet: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// generate token
		token, err := GenerateToken(u.ID, u.Role)
		if err != nil {
			http.Error(w, "Failed to generate token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resp := SignupResponse{Token: token}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
}

// LoginHandler handles POST /auth/login
func LoginHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
			return
		}

		userRepo := user.NewRepository(db)

		u, err := userRepo.FindByEmail(req.Email)
		if err != nil {
			http.Error(w, "User lookup failed: "+err.Error(), http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
			http.Error(w, "Invalid credentials: "+err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := GenerateToken(u.ID, u.Role)
		if err != nil {
			http.Error(w, "Failed to generate token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resp := LoginResponse{Token: token}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
}
