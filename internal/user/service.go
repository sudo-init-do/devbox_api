package user

import "database/sql"

type Service struct {
	repo *Repository
}

func NewService(db *sql.DB) *Service {
	return &Service{repo: NewRepository(db)}
}

func (s *Service) Create(email, password, role string) (*User, error) {
	return s.repo.Create(email, password, role)
}

func (s *Service) FindByEmail(email string) (*User, error) {
	return s.repo.FindByEmail(email)
}
