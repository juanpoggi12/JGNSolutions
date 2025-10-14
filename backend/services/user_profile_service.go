package services

import (
	"errors"
	"strings"
	"time"

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
	existing, _ := service.repository.ObtenerPerfilPorUserID(actor.UserID)
	if existing.UserID == actor.UserID {
		return models.UserProfile{}, errors.New("el perfil ya existe")
	}

	// Convertir DTO → modelo (con el userID del actor)
	perfil, err := utils.ConvertUserProfileCreateRequestToModel(req, actor.UserID)
	if err != nil {
		return models.UserProfile{}, err
	}

	_, err = service.repository.InsertarPerfil(perfil)
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
	perfil, err := service.repository.ObtenerPerfilPorUserID(actor.UserID)
	if err != nil {
		return models.UserProfile{}, errors.New("perfil no encontrado")
	}
	return perfil, nil
}

// Actualizar perfil (usuario propio o admin)
func (service *UserProfileService) UpdateProfile(actor Actor, req dto.UserProfileUpdateRequest) (models.UserProfile, error) {
	if actor.UserID.IsZero() {
		return models.UserProfile{}, errors.New("ID de usuario inválido")
	}

	// Obtener perfil actual
	perfilActual, err := service.repository.ObtenerPerfilPorUserID(actor.UserID)
	if err != nil {
		return models.UserProfile{}, errors.New("perfil no encontrado")
	}

	// Validar permisos
	if strings.ToLower(actor.Role) != "admin" && perfilActual.UserID != actor.UserID {
		return models.UserProfile{}, errors.New("no tienes permiso para modificar este perfil")
	}

	// Aplicar solo los campos presentes en el DTO
	if req.FullName != nil {
		perfilActual.FullName = *req.FullName
	}
	if req.BirthDate != nil {
		fecha, err := time.Parse("2006-01-02", *req.BirthDate)
		if err == nil {
			perfilActual.BirthDate = fecha
		}
	}
	if req.WeightKg != nil {
		perfilActual.WeightKg = *req.WeightKg
	}
	if req.HeightCm != nil {
		perfilActual.HeightCm = *req.HeightCm
	}
	if req.Level != nil {
		perfilActual.Level = models.Nivel(*req.Level)
	}
	if req.Goal != nil {
		perfilActual.Goal = models.Objetivo(*req.Goal)
	}

	perfilActual.UpdatedAt = time.Now()

	_, err = service.repository.ModificarPerfil(perfilActual)
	if err != nil {
		return models.UserProfile{}, err
	}

	if service.logService != nil {
		service.logService.RecordAction(actor.UserID, "perfil:modificado "+perfilActual.UserID.Hex())
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
	_, err = service.repository.EliminarPerfil(objectID)
	if err != nil {
		return err
	}

	if service.logService != nil {
		service.logService.RecordAction(actor.UserID, "perfil:eliminado "+id)
	}

	return nil
}
