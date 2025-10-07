package utils

import (
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
)

// CreateRequest → Model
// Ya no toma el UserID del DTO, se asigna externamente en el service usando actor.UserID
func ConvertUserProfileCreateRequestToModel(req dto.UserProfileCreateRequest) (models.UserProfile, error) {
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return models.UserProfile{}, err
	}

	return models.UserProfile{
		// UserID se asignará en el service
		FullName:  req.FullName,
		BirthDate: birthDate,
		WeightKg:  req.WeightKg,
		HeightCm:  req.HeightCm,
		Level:     models.Nivel(req.Level),
		Goal:      models.Objetivo(req.Goal),
		UpdatedAt: time.Now(),
	}, nil
}

// UpdateRequest → aplica cambios sobre el modelo
func ApplyUserProfileUpdateToModel(up *models.UserProfile, req dto.UserProfileUpdateRequest) error {
	if req.FullName != nil {
		up.FullName = *req.FullName
	}
	if req.BirthDate != nil {
		bd, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			return err
		}
		up.BirthDate = bd
	}
	if req.WeightKg != nil {
		up.WeightKg = *req.WeightKg
	}
	if req.HeightCm != nil {
		up.HeightCm = *req.HeightCm
	}
	if req.Level != nil {
		up.Level = models.Nivel(*req.Level)
	}
	if req.Goal != nil {
		up.Goal = models.Objetivo(*req.Goal)
	}
	up.UpdatedAt = time.Now()
	return nil
}

// Model → Response
// Quita UserID si el DTO ya no lo tiene
func ConvertUserProfileModelToResponse(up models.UserProfile) dto.UserProfileResponse {
	return dto.UserProfileResponse{
		ID:        up.ID.Hex(),
		FullName:  up.FullName,
		BirthDate: up.BirthDate.Format("2006-01-02"),
		WeightKg:  up.WeightKg,
		HeightCm:  up.HeightCm,
		Level:     string(up.Level),
		Goal:      string(up.Goal),
		UpdatedAt: up.UpdatedAt.Format(time.RFC3339),
	}
}

// SearchRequest → filtro Mongo
// Ya no usa UserID del body; el filtro por usuario lo manejará el service usando actor.UserID
func BuildUserProfileSearchFilter(search dto.UserProfileSearchRequest, userID string) bson.M {
	filter := bson.M{}

	if userID != "" {
		filter["userId"] = ToObjectID(userID)
	}
	if search.Name != "" {
		filter["fullName"] = bson.M{"$regex": search.Name, "$options": "i"}
	}
	if search.Level != "" {
		filter["level"] = search.Level
	}
	if search.Goal != "" {
		filter["goal"] = search.Goal
	}
	return filter
}
