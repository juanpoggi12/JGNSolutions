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

type UserServiceInterface interface {
	GetUserByID(actor Actor, id string) (models.User, error)
	DeleteUser(actor Actor, id string) error
}

type UserService struct {
	repository repositories.UserRepositoryInterface
	logService *LogService
}

func NewUserService(repository repositories.UserRepositoryInterface, logService *LogService) *UserService {
	return &UserService{repository: repository, logService: logService}
}

func (service *UserService) CreateUser(actor Actor, req dto.UserCreateRequest) (models.User, error) {
	// Solo los administradores pueden crear nuevos usuarios
	if strings.ToLower(actor.Role) != "admin" {
		return models.User{}, errors.New("solo los administradores pueden crear usuarios")
	}

	// Validar que venga contraseña
	if req.Password == "" {
		return models.User{}, errors.New("la contraseña es obligatoria")
	}

	// Convertir DTO → Modelo (ya hashea la contraseña)
	user, err := utils.ConvertUserCreateRequestToModel(req)
	if err != nil {
		return models.User{}, errors.New("error al crear modelo de usuario")
	}

	// Rol por defecto: user
	if user.Role == "" {
		user.Role = "user"
	}

	// Insertar en la base
	resultado, err := service.repository.InsertUser(user)
	if err != nil {
		return models.User{}, err
	}

	// Asignar el ID generado
	if oid, ok := resultado.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
		if service.logService != nil {
			service.logService.RecordAction(actor.UserID, "usuario:creado "+user.Email)
		}
		return user, nil
	}

	return models.User{}, errors.New("error al crear usuario")
}

func (service *UserService) GetUserByID(actor Actor, id string) (models.User, error) {
	usuario, err := service.repository.GetUserByID(id)
	if err != nil {
		return models.User{}, errors.New("usuario no encontrado")
	}

	// Un usuario solo puede ver su propio perfil, excepto que sea admin
	if actor.Role != "admin" && usuario.ID != actor.UserID {
		return models.User{}, errors.New("no tienes permiso para acceder a este usuario")
	}

	return usuario, nil
}

func (service *UserService) DeleteUser(actor Actor, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID inválido")
	}

	// Solo admin puede eliminar usuarios
	if actor.Role != "admin" {
		return errors.New("solo los administradores pueden eliminar usuarios")
	}

	resultado, err := service.repository.DeleteUser(objectID)
	if err != nil {
		return err
	}
	if resultado.DeletedCount == 0 {
		return errors.New("usuario no encontrado")
	}

	if service.logService != nil {
		service.logService.RecordAction(actor.UserID, "usuario:eliminado "+id)
	}

	return nil
}
