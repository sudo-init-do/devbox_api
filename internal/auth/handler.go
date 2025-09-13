package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/sudo-init-do/devbox_api/internal/user"
	"golang.org/x/crypto/bcrypt"
)

// DTOs
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

// Handler struct
type Handler struct {
	userService *user.Service
}

func NewHandler(us *user.Service) *Handler {
	return &Handler{userService: us}
}

// Signup handles user registration
func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	// create user
	u := user.User{
		Email:    req.Email,
		Password: string(hashed),
		Role:     "user",
	}

	if err := h.userService.Create(&u); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	token, _ := GenerateToken(u.ID, u.Role)
	json.NewEncoder(w).Encode(SignupResponse{Token: token})
}

// Login handles user login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// find user
	u, err := h.userService.FindByEmail(req.Email)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	// compare password
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, _ := GenerateToken(u.ID, u.Role)
	json.NewEncoder(w).Encode(SignupResponse{Token: token})
}
