package user

import (
	"database/sql"

	"github.com/google/uuid"
)

// Repository handles DB operations for users
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a new user into the DB
func (r *Repository) Create(email, passwordHash, role string) (*User, error) {
	id := uuid.New().String()

	_, err := r.db.Exec(
		`INSERT INTO users (id, email, password_hash, role) VALUES ($1, $2, $3, $4)`,
		id, email, passwordHash, role,
	)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
	}, nil
}

// FindByEmail retrieves a user by email
func (r *Repository) FindByEmail(email string) (*User, error) {
	row := r.db.QueryRow(`SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1`, email)

	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}
