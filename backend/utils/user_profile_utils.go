package utils

import (
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateRequest → Model
func ConvertUserProfileCreateRequestToModel(req dto.UserProfileCreateRequest) (models.UserProfile, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return models.UserProfile{}, err
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return models.UserProfile{}, err
	}

	return models.UserProfile{
		UserID:    userID,
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
func ConvertUserProfileModelToResponse(up models.UserProfile) dto.UserProfileResponse {
	return dto.UserProfileResponse{
		ID:        up.ID.Hex(),
		UserID:    up.UserID.Hex(),
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
func BuildUserProfileSearchFilter(search dto.UserProfileSearchRequest) bson.M {
	filter := bson.M{}
	if search.UserID != "" {
		filter["userId"] = ToObjectID(search.UserID)
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
