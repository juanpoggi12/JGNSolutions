package services

import (
	"errors"

	"github.com/juanpoggi12/JGNSolutions/backend/repositories"
)

type AdminServiceInterface interface {
	CountUsers(actor Actor) (int64, error)
	CountExercises(actor Actor) (int64, error)
	CountRoutines(actor Actor) (int64, error)
	CountWorkoutSessions(actor Actor) (int64, error)
}

type AdminService struct {
	adminRepository repositories.AdminRepositoryInterface
}

func NewAdminService(repository repositories.AdminRepositoryInterface) *AdminService {
	return &AdminService{adminRepository: repository}
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
