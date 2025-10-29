package utils

import (
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"
)

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
		CreatedAt: up.CreatedAt.Format(time.RFC3339),
		UpdatedAt: up.UpdatedAt.Format(time.RFC3339),
	}
}
