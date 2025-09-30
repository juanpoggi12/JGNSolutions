package utils

import (
	"JGNSolutions/backend/dto"
	"JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Request -> Model
func ConvertUserRequestToModel(req dto.UserRequest) models.User {
	return models.User{
		Email:    req.Email,
		Password: req.Password, // ojo: en el service deberÃ­as hashearla antes de guardar
		Role:     req.Role,
	}
}

// Model -> Response
func ConvertUserModelToResponse(u models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:    u.ID.Hex(),
		Email: u.Email,
		Role:  u.Role,
	}
}

// Helper para convertir string a ObjectID
func ToObjectID(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID
	}
	return oid
}
