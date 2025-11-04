package services

import (
	"errors"
	"strings"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"
	"github.com/juanpoggi12/JGNSolutions/backend/utils"
)

type UserProfileServiceInterface interface {
	GetProfile(actor Actor) (models.UserProfile, error)
	UpdateProfile(actor Actor, req dto.UserProfileUpdateRequest) (models.UserProfile, error)
}

type UserProfileService struct {
	repository repositories.UserProfileRepositoryInterface
	logService *LogService
}

func NewUserProfileService(repo repositories.UserProfileRepositoryInterface, logService *LogService) *UserProfileService {
	return &UserProfileService{repository: repo, logService: logService}
}

func (service *UserProfileService) GetProfile(actor Actor) (models.UserProfile, error) {
	if actor.UserID.IsZero() {
		return models.UserProfile{}, errors.New("ID de usuario inválido")
	}
	perfil, err := service.repository.GetProfileByUserID(actor.UserID)
	if err != nil {
		return models.UserProfile{}, errors.New("perfil no encontrado")
	}
	return perfil, nil
}

func (service *UserProfileService) UpdateProfile(actor Actor, id string, req dto.UserProfileUpdateRequest) (models.UserProfile, error) {
	if actor.UserID.IsZero() {
		return models.UserProfile{}, errors.New("ID de usuario inválido")
	}

	perfilActual, err := service.repository.GetProfileByID(id)
	if err != nil {
		return models.UserProfile{}, errors.New("perfil no encontrado")
	}

	if strings.ToLower(actor.Role) != "admin" && perfilActual.UserID != actor.UserID {
		return models.UserProfile{}, errors.New("no tienes permiso para modificar este perfil")
	}

	if err := utils.ApplyUserProfileUpdateToModel(&perfilActual, req); err != nil {
		return models.UserProfile{}, err
	}

	_, err = service.repository.UpdateProfile(perfilActual)
	if err != nil {
		return models.UserProfile{}, err
	}

	if service.logService != nil {
		service.logService.RecordAction(actor.UserID, "perfil actualizado: "+perfilActual.ID.Hex())
	}

	return perfilActual, nil
}
