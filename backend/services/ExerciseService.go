package services

import (
	"JGNSolutions/backend/models"
	"JGNSolutions/backend/repositories"
	"context"
)

type ExerciseService struct {
	repo *repositories.ExerciseRepository
}

func NewExerciseService(r *repositories.ExerciseRepository) *ExerciseService {
	return &ExerciseService{repo: r}
}

func (s *ExerciseService) Create(ctx context.Context, exercise *models.Exercise) error {
	return s.repo.Create(ctx, exercise)
}

func (s *ExerciseService) GetByID(ctx context.Context, id string) (*models.Exercise, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ExerciseService) Update(ctx context.Context, exercise *models.Exercise) error {
	return s.repo.Update(ctx, exercise)
}

func (s *ExerciseService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
