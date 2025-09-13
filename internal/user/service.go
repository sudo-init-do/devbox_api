package user

import (
	"database/sql"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(u *User) error {
	// generate ID in DB
	_, err := s.db.Exec(
		"INSERT INTO users (id, email, password, role) VALUES (gen_random_uuid(), $1, $2, $3)",
		u.Email, u.Password, u.Role,
	)
	return err
}

func (s *Service) FindByEmail(email string) (*User, error) {
	row := s.db.QueryRow("SELECT id, email, password, role FROM users WHERE email=$1", email)

	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Role); err != nil {
		return nil, err
	}
	return &u, nil
}
