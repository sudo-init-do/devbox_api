package user

import (
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a new user
func (r *Repository) Create(email, password, role string) (*User, error) {
	var id string
	err := r.db.QueryRow(`
		INSERT INTO users (email, password, role, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`, email, password, role).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &User{
		ID:        id,
		Email:     email,
		Password:  password,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// FindByEmail retrieves a user by email
func (r *Repository) FindByEmail(email string) (*User, error) {
	row := r.db.QueryRow(`
		SELECT id, email, password, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email)

	var u User
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}
