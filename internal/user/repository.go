package user

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(email, password, role string) (*User, error) {
	id := uuid.New()
	_, err := r.db.Exec(`INSERT INTO users (id, email, password_hash, role, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		id, email, password, role, time.Now(), time.Now())
	if err != nil {
		return nil, err
	}

	return &User{ID: id, Email: email, Password: password, Role: role}, nil
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	row := r.db.QueryRow(`SELECT id, email, password_hash, role, created_at, updated_at 
		FROM users WHERE email = $1`, email)

	var u User
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
