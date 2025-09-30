package services

import (
	"JGNSolutions/backend/repositories"
	"context"
)

type AdminService struct {
	userRepo     *repositories.UserRepository
	exerciseRepo *repositories.ExerciseRepository
	routineRepo  *repositories.RoutineRepository
	workoutRepo  *repositories.WorkoutRepository
}

func NewAdminService(u *repositories.UserRepository, e *repositories.ExerciseRepository, r *repositories.RoutineRepository, w *repositories.WorkoutRepository) *AdminService {
	return &AdminService{
		userRepo:     u,
		exerciseRepo: e,
		routineRepo:  r,
		workoutRepo:  w,
	}
}

// Métodos básicos (placeholders, luego los completamos con queries)

func (s *AdminService) GetTotalUsers(ctx context.Context) (int, error) {
	return s.userRepo.Count(ctx)
}

func (s *AdminService) GetMostPopularExercises(ctx context.Context) ([]string, error) {
	return s.exerciseRepo.MostPopular(ctx)
}

func (s *AdminService) GetMostUsedRoutines(ctx context.Context) ([]string, error) {
	return s.routineRepo.MostUsed(ctx)
}
