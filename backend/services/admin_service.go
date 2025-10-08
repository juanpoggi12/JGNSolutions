package services

import (
	"errors"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"
)

type AdminServiceInterface interface {
	CountUsers(actor Actor) (int64, error)
	CountExercises(actor Actor) (int64, error)
	CountRoutines(actor Actor) (int64, error)
	CountWorkoutSessions(actor Actor) (int64, error)
	ListUsers(actor Actor) ([]dto.UserListResponse, error)
	TopExercises(actor Actor, limit int) ([]dto.ExerciseStatResponse, error)
	TopRoutines(actor Actor, limit int) ([]dto.RoutineStatResponse, error)
}

type AdminService struct {
	adminRepository repositories.AdminRepositoryInterface
	userRepository  repositories.UserRepositoryInterface
}

func NewAdminService(adminRepo repositories.AdminRepositoryInterface, userRepo repositories.UserRepositoryInterface) *AdminService {
	return &AdminService{adminRepository: adminRepo, userRepository: userRepo}
}

func (service *AdminService) CountUsers(actor Actor) (int64, error) {
	if actor.Role != "admin" {
		return 0, errors.New("no tienes permiso para acceder a estas estadísticas")
	}
	return service.adminRepository.ContarDocumentos("users")
}

func (service *AdminService) CountExercises(actor Actor) (int64, error) {
	if actor.Role != "admin" {
		return 0, errors.New("no tienes permiso para acceder a estas estadísticas")
	}
	return service.adminRepository.ContarDocumentos("exercises")
}

func (service *AdminService) CountRoutines(actor Actor) (int64, error) {
	if actor.Role != "admin" {
		return 0, errors.New("no tienes permiso para acceder a estas estadísticas")
	}
	return service.adminRepository.ContarDocumentos("routines")
}

func (service *AdminService) CountWorkoutSessions(actor Actor) (int64, error) {
	if actor.Role != "admin" {
		return 0, errors.New("no tienes permiso para acceder a estas estadísticas")
	}
	return service.adminRepository.ContarDocumentos("workout_sessions")
}

func (s *AdminService) ListUsers(actor Actor) ([]dto.UserListResponse, error) {
	if actor.Role != "admin" {
		return nil, errors.New("no tienes permiso para ver usuarios")
	}
	users, err := s.userRepository.ListarUsuariosBasico()
	if err != nil {
		return nil, err
	}

	out := make([]dto.UserListResponse, 0, len(users))
	for _, u := range users {
		out = append(out, dto.UserListResponse{
			ID:    u.ID.Hex(),
			Email: u.Email,
			Role:  string(u.Role), // 👈 conversión explícita
		})
	}
	return out, nil
}

func (s *AdminService) TopExercises(actor Actor, limit int) ([]dto.ExerciseStatResponse, error) {
	if actor.Role != "admin" {
		return nil, errors.New("no tienes permiso para ver estadísticas")
	}
	rows, err := s.adminRepository.TopExercisesByEntries(limit)
	if err != nil {
		return nil, err
	}

	out := make([]dto.ExerciseStatResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, dto.ExerciseStatResponse{
			ExerciseID: r.ID.Hex(),
			Name:       r.Name,
			UsageCount: r.UsageCount,
		})
	}
	return out, nil
}

func (s *AdminService) TopRoutines(actor Actor, limit int) ([]dto.RoutineStatResponse, error) {
	if actor.Role != "admin" {
		return nil, errors.New("no tienes permiso para ver estadísticas")
	}
	rows, err := s.adminRepository.TopRoutinesBySessions(limit)
	if err != nil {
		return nil, err
	}

	out := make([]dto.RoutineStatResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, dto.RoutineStatResponse{
			RoutineID:  r.ID.Hex(),
			Name:       r.Name,
			UsageCount: r.UsageCount,
		})
	}
	return out, nil
}
