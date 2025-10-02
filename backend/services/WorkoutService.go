package services

import (
    "github.com/juanpoggi12/JGNSolutions/backend/models"
    "github.com/juanpoggi12/JGNSolutions/backend/repositories"
	"context"
)

type WorkoutService struct {
	repo *repositories.WorkoutRepository
}

func NewWorkoutService(r *repositories.WorkoutRepository) *WorkoutService {
	return &WorkoutService{repo: r}
}

func (s *WorkoutService) Create(ctx context.Context, workout *models.Workout) error {
	return s.repo.Create(ctx, workout)
}

func (s *WorkoutService) GetByID(ctx context.Context, id string) (*models.Workout, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *WorkoutService) Update(ctx context.Context, workout *models.Workout) error {
	return s.repo.Update(ctx, workout)
}

func (s *WorkoutService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
