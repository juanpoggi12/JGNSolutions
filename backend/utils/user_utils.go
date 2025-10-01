package utils

import (
	"time"

	"github.com/juanpoggi12/JGNSolutions/dto"
	"github.com/juanpoggi12/JGNSolutions/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ConvertUserCreateRequestToModel converts DTO to model for creating a user
func ConvertUserCreateRequestToModel(req dto.UserCreateRequest) models.User {
	return models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password, // service should hash before saving
		Role:         models.Role(req.Role),
		IsActive:     req.IsActive != nil && *req.IsActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// ConvertUserRequestToModel converts general DTO to model (used in some handlers)
func ConvertUserRequestToModel(req dto.UserUpdateRequest) models.User {
	// partial update mapping; caller should apply non-nil fields
	u := models.User{}
	if req.Email != nil {
		u.Email = *req.Email
	}
	if req.Password != nil {
		u.PasswordHash = *req.Password
	}
	if req.Role != nil {
		u.Role = models.Role(*req.Role)
	}
	return u
}

// ConvertUserModelToResponse maps model -> DTO response
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

// ToObjectID converts hex string to primitive.ObjectID, returns NilObjectID on error
func ToObjectID(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID
	}
	return oid
}
