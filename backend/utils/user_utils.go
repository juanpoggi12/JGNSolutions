package utils

import (
	"strings"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// CreateRequest → Model (con hash de password)
func ConvertUserCreateRequestToModel(req dto.UserCreateRequest) (models.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	// 🔹 Normalizar el rol a minúsculas (evita errores con ADMIN vs admin)
	role := strings.ToLower(strings.TrimSpace(req.Role))
	if role == "" {
		role = "user" // valor por defecto si no se envía
	}

	return models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashed),
		Role:         models.Role(role),
		IsActive:     req.IsActive != nil && *req.IsActive, // default false si no viene
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

// Model → Response
func ConvertUserModelToResponse(u models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID.Hex(),
		Username:  u.Username,
		Email:     u.Email,
		Role:      string(u.Role),
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
		IsActive:  u.IsActive,
	}
}

// ToObjectID helper (ya lo tenías)
func ToObjectID(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID
	}
	return oid
}
