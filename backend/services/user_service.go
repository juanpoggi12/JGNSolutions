package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"
	"github.com/juanpoggi12/JGNSolutions/backend/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserServiceInterface interface {
	CreateUser(actor Actor, user models.User) (models.User, error)
	GetUserByID(actor Actor, id string) (models.User, error)
	UpdateUser(actor Actor, id string, user models.User) (models.User, error)
	DeleteUser(actor Actor, id string) error
	ChangePassword(ctx context.Context, userID primitive.ObjectID, req dto.ChangePasswordRequest) error
}

type UserService struct {
	repository repositories.UserRepositoryInterface
}

func NewUserService(repository repositories.UserRepositoryInterface) *UserService {
	return &UserService{repository: repository}
}

func (service *UserService) CreateUser(actor Actor, user models.User) (models.User, error) {
	// Solo los administradores pueden crear nuevos usuarios
	if actor.Role != "admin" {
		return models.User{}, errors.New("solo los administradores pueden crear usuarios")
	}

	// Validar que venga una contraseña
	if user.PasswordHash == "" {
		return models.User{}, errors.New("la contraseña es obligatoria")
	}

	// Hashear la contraseña antes de guardar
	hash, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		return models.User{}, errors.New("error al hashear la contraseña")
	}
	user.PasswordHash = hash

	// Normalizar el rol a minúsculas (por coherencia con el resto del código)
	user.Role = models.Role(strings.ToLower(string(user.Role)))

	// Inicializar metadatos
	user.ID = primitive.NilObjectID
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Insertar en la base de datos
	resultado, err := service.repository.InsertarUsuario(user)
	if err != nil {
		return models.User{}, err
	}

	// Si la inserción fue correcta, asignar el ObjectID generado
	if oid, ok := resultado.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
		return user, nil
	}

	return models.User{}, errors.New("error al crear usuario")
}

func (service *UserService) GetUserByID(actor Actor, id string) (models.User, error) {
	usuario, err := service.repository.ObtenerUsuarioPorID(id)
	if err != nil {
		return models.User{}, errors.New("usuario no encontrado")
	}

	// Un usuario solo puede ver su propio perfil, excepto que sea admin
	if actor.Role != "admin" && usuario.ID != actor.UserID {
		return models.User{}, errors.New("no tienes permiso para acceder a este usuario")
	}

	return usuario, nil
}

func (service *UserService) UpdateUser(actor Actor, id string, req dto.UserUpdateRequest) (models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, errors.New("ID inválido")
	}

	userDB, err := service.repository.ObtenerUsuarioPorID(id)
	if err != nil {
		return models.User{}, errors.New("usuario no encontrado")
	}

	// Solo admin o el mismo usuario puede actualizar su perfil
	if actor.Role != "admin" && actor.UserID != objectID {
		return models.User{}, errors.New("no tienes permiso para modificar este usuario")
	}

	// Aplicar solo los campos presentes
	if req.Username != nil {
		userDB.Username = *req.Username
	}
	if req.Email != nil {
		userDB.Email = *req.Email
	}
	if req.Password != nil && *req.Password != "" {
		hash, err := utils.HashPassword(*req.Password)
		if err != nil {
			return models.User{}, err
		}
		userDB.PasswordHash = hash
	}
	if req.Role != nil && *req.Role != "" {
		userDB.Role = models.Role(strings.ToLower(*req.Role))
	}

	userDB.UpdatedAt = time.Now()

	_, err = service.repository.ModificarUsuario(userDB)
	if err != nil {
		return models.User{}, err
	}

	return userDB, nil
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

	resultado, err := service.repository.EliminarUsuario(objectID)
	if err != nil {
		return err
	}
	if resultado.DeletedCount == 0 {
		return errors.New("usuario no encontrado")
	}

	return nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID primitive.ObjectID, req dto.ChangePasswordRequest) error {
	// Buscar el usuario en la base de datos
	user, err := s.repository.ObtenerUsuarioPorID(userID.Hex()) // 🔹 Convertimos ObjectID a string
	if err != nil {
		return err
	}

	// Verificar la contraseña actual
	if !utils.VerifyPassword(req.OldPassword, user.PasswordHash) {
		return errors.New("contraseña actual incorrecta")
	}

	// Hashear la nueva contraseña
	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Guardar cambios en la base (ignoramos el resultado del Update)
	_, err = s.repository.ActualizarPassword(userID, hashed)
	return err
}
