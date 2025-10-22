package services

import (
	"errors"
	"strings"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"
	"github.com/juanpoggi12/JGNSolutions/backend/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 🧩 Interface del servicio
type UserProfileServiceInterface interface {
	CreateProfile(actor Actor, req dto.UserProfileCreateRequest) (models.UserProfile, error)
	GetProfile(actor Actor) (models.UserProfile, error)
	UpdateProfile(actor Actor, req dto.UserProfileUpdateRequest) (models.UserProfile, error)
	DeleteProfile(actor Actor, id string) error
}

// 🧩 Implementación del servicio
type UserProfileService struct {
	repository repositories.UserProfileRepositoryInterface
	logService *LogService
}

func NewUserProfileService(repo repositories.UserProfileRepositoryInterface, logService *LogService) *UserProfileService {
	return &UserProfileService{repository: repo, logService: logService}
}

// Crear perfil nuevo (solo si el usuario no lo tiene aún)
func (service *UserProfileService) CreateProfile(actor Actor, req dto.UserProfileCreateRequest) (models.UserProfile, error) {
	if actor.UserID.IsZero() {
		return models.UserProfile{}, errors.New("ID de usuario inválido")
	}

	// Evitar duplicados
	existing, _ := service.repository.GetProfileByUserID(actor.UserID)
	if existing.UserID == actor.UserID {
		return models.UserProfile{}, errors.New("el perfil ya existe")
	}

	// Convertir DTO → modelo (con el userID del actor)
	perfil, err := utils.ConvertUserProfileCreateRequestToModel(req, actor.UserID)
	if err != nil {
		return models.UserProfile{}, err
	}

	_, err = service.repository.InsertProfile(perfil)
	if err != nil {
		return models.UserProfile{}, err
	}

	if service.logService != nil {
		service.logService.RecordAction(actor.UserID, "perfil:creado "+actor.UserID.Hex())
	}

	return perfil, nil
}

// Obtener el perfil del usuario autenticado
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

	// Buscar el perfil objetivo
	perfilActual, err := service.repository.GetProfileByID(id)
	if err != nil {
		return models.UserProfile{}, errors.New("perfil no encontrado")
	}

	// Permitir modificar si:
	// - el actor es admin, o
	// - el actor es dueño del perfil
	if strings.ToLower(actor.Role) != "admin" && perfilActual.UserID != actor.UserID {
		return models.UserProfile{}, errors.New("no tienes permiso para modificar este perfil")
	}

	// Aplicar cambios usando helper del utils
	if err := utils.ApplyUserProfileUpdateToModel(&perfilActual, req); err != nil {
		return models.UserProfile{}, err
	}

	// Guardar cambios
	_, err = service.repository.UpdateProfile(perfilActual)
	if err != nil {
		return models.UserProfile{}, err
	}

	// Log opcional
	if service.logService != nil {
		service.logService.RecordAction(actor.UserID, "perfil actualizado: "+perfilActual.ID.Hex())
	}

	return perfilActual, nil
}


// Eliminar un perfil (solo admin)
func (service *UserProfileService) DeleteProfile(actor Actor, id string) error {
	if strings.ToLower(actor.Role) != "admin" {
		return errors.New("solo los administradores pueden eliminar perfiles")
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID inválido")
	}
	_, err = service.repository.DeleteProfile(objectID)
	if err != nil {
		return err
	}

	if service.logService != nil {
		service.logService.RecordAction(actor.UserID, "perfil:eliminado "+id)
	}

	return nil
}
