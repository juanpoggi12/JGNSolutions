package services

import (
	"errors"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserProfileServiceInterface interface {
	CreateProfile(actor Actor, perfil models.UserProfile) (models.UserProfile, error)
	GetProfile(actor Actor) (models.UserProfile, error)
	UpdateProfile(actor Actor, perfil models.UserProfile) (models.UserProfile, error)
	DeleteProfile(actor Actor, id string) error
}

type UserProfileService struct {
	repository repositories.UserProfileRepositoryInterface
}

func NewUserProfileService(repo repositories.UserProfileRepositoryInterface) *UserProfileService {
	return &UserProfileService{repository: repo}
}

func (service *UserProfileService) CreateProfile(actor Actor, perfil models.UserProfile) (models.UserProfile, error) {
	if actor.UserID.IsZero() {
		return models.UserProfile{}, errors.New("ID de usuario inválido")
	}
	perfil.UserID = actor.UserID
	perfil.UpdatedAt = time.Now()
	_, err := service.repository.InsertarPerfil(perfil)
	return perfil, err
}

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

func (service *UserProfileService) UpdateProfile(actor Actor, perfil models.UserProfile) (models.UserProfile, error) {
	if actor.UserID.IsZero() {
		return models.UserProfile{}, errors.New("ID de usuario inválido")
	}
	perfil.UserID = actor.UserID
	perfil.UpdatedAt = time.Now()
	_, err := service.repository.ModificarPerfil(perfil)
	return perfil, err
}

func (service *UserProfileService) DeleteProfile(actor Actor, id string) error {
	if actor.Role != "admin" {
		return errors.New("solo los administradores pueden eliminar perfiles")
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID inválido")
	}
	_, err = service.repository.EliminarPerfil(objectID)
	return err
}
