package utils

import (
	"time"

    "github.com/juanpoggi12/JGNSolutions/backend/dto"
    "github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// CreateRequest → Model (con hash de password)
func ConvertUserCreateRequestToModel(req dto.UserCreateRequest) (models.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashed),
		Role:         models.Role(req.Role),
		IsActive:     req.IsActive != nil && *req.IsActive, // default false si no viene
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

// UpdateRequest → aplica cambios sobre el modelo existente
func ApplyUserUpdateToModel(u *models.User, req dto.UserUpdateRequest) error {
	if req.Username != nil {
		u.Username = *req.Username
	}
	if req.Email != nil {
		u.Email = *req.Email
	}
	if req.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.PasswordHash = string(hashed)
	}
	if req.Role != nil {
		u.Role = models.Role(*req.Role)
	}
	if req.IsActive != nil {
		u.IsActive = *req.IsActive
	}
	u.UpdatedAt = time.Now()
	return nil
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

// SearchRequest → filtro Mongo
func BuildUserSearchFilter(search dto.UserSearchRequest) bson.M {
	filter := bson.M{}
	if search.Username != "" {
		filter["username"] = bson.M{"$regex": search.Username, "$options": "i"}
	}
	if search.Email != "" {
		filter["email"] = bson.M{"$regex": search.Email, "$options": "i"}
	}
	if search.Role != "" {
		filter["role"] = search.Role
	}
	if search.IsActive != nil {
		filter["isActive"] = *search.IsActive
	}
	return filter
}

// ToObjectID helper (ya lo tenías)
func ToObjectID(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID
	}
	return oid
}
