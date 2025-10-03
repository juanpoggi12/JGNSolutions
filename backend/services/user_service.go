package services

import (
	"errors"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserServiceInterface interface {
	CreateUser(user models.User) (models.User, error)
	GetUserByID(id string) (models.User, error)
	UpdateUser(id string, user models.User) (models.User, error)
	DeleteUser(id string) error
}

type UserService struct {
	repository repositories.UserRepositoryInterface
}

func NewUserService(repository repositories.UserRepositoryInterface) *UserService {
	return &UserService{repository: repository}
}

func (service *UserService) CreateUser(user models.User) (models.User, error) {
	user.ID = primitive.NilObjectID
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	resultado, err := service.repository.InsertarUsuario(user)
	if err != nil {
		return models.User{}, err
	}

	if oid, ok := resultado.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
		return user, nil
	}

	return models.User{}, errors.New("error al crear usuario")
}

func (service *UserService) GetUserByID(id string) (models.User, error) {
	usuario, err := service.repository.ObtenerUsuarioPorID(id)
	if err != nil {
		return models.User{}, errors.New("usuario no encontrado")
	}
	return usuario, nil
}

func (service *UserService) UpdateUser(id string, user models.User) (models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, errors.New("ID inválido")
	}

	user.ID = objectID
	user.UpdatedAt = time.Now()
	_, err = service.repository.ModificarUsuario(user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (service *UserService) DeleteUser(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID inválido")
	}

	resultado, err := service.repository.EliminarUsuario(objectID)
	if err != nil {
		return err
	}
	if resultado.DeletedCount == 0 {
		return errors.New("usuario no encontrado")
	}

	return nil
}
