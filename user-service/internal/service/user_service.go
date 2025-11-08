package service

import (
    "user-service/internal/repository"
)

type UserService struct {
    repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
    return &UserService{repo: r}
}

func (s *UserService) Create(u *repository.User) error {
    if u.Role != "PATIENT" && u.Role != "DOCTOR" {
        return &ErrInvalidRole{}
    }
    return s.repo.Create(u)
}

func (s *UserService) GetAll() ([]*repository.User, error) {
    return s.repo.GetAll()
}

func (s *UserService) GetByID(id string) (*repository.User, error) {
    return s.repo.GetByID(id)
}

func (s *UserService) Update(id string, u *repository.User) error {
    return s.repo.Update(id, u)
}

func (s *UserService) Delete(id string) error {
    return s.repo.Delete(id)
}

type ErrInvalidRole struct{}

func (e *ErrInvalidRole) Error() string { return "invalid role, must be PATIENT or DOCTOR" }
