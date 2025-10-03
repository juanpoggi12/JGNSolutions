package services

import (
    "github.com/juanpoggi12/JGNSolutions/backend/models"
    "github.com/juanpoggi12/JGNSolutions/backend/repositories"
	"context"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(r *repositories.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) Create(ctx context.Context, user *models.User) error {
	return s.repo.Create(ctx, user)
}

func (s *UserService) GetByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, user *models.User) error {
	return s.repo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
