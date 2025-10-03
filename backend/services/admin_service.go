package services

import "github.com/juanpoggi12/JGNSolutions/backend/repositories"

type AdminServiceInterface interface {
	CountUsers() (int64, error)
	CountExercises() (int64, error)
	CountRoutines() (int64, error)
	CountWorkoutSessions() (int64, error)
}

type AdminService struct {
	adminRepository repositories.AdminRepositoryInterface
}

func NewAdminService(repository repositories.AdminRepositoryInterface) *AdminService {
	return &AdminService{adminRepository: repository}
}

func (service *AdminService) CountUsers() (int64, error) {
	return service.adminRepository.ContarDocumentos("users")
}

func (service *AdminService) CountExercises() (int64, error) {
	return service.adminRepository.ContarDocumentos("exercises")
}

func (service *AdminService) CountRoutines() (int64, error) {
	return service.adminRepository.ContarDocumentos("routines")
}

func (service *AdminService) CountWorkoutSessions() (int64, error) {
	return service.adminRepository.ContarDocumentos("workout_sessions")
}
