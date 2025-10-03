package services

import (
    "github.com/juanpoggi12/JGNSolutions/backend/models"
    "github.com/juanpoggi12/JGNSolutions/backend/repositories"
	"context"
)

type RoutineService struct {
	repo *repositories.RoutineRepository
}

func NewRoutineService(r *repositories.RoutineRepository) *RoutineService {
	return &RoutineService{repo: r}
}

func (s *RoutineService) Create(ctx context.Context, routine *models.Routine) error {
	return s.repo.Create(ctx, routine)
}

func (s *RoutineService) GetByID(ctx context.Context, id string) (*models.Routine, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *RoutineService) Update(ctx context.Context, routine *models.Routine) error {
	return s.repo.Update(ctx, routine)
}

func (s *RoutineService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
