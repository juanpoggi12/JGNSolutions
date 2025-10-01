package utils

import (
	"github.com/juanpoggi12/JGNSolutions/dto"
	"github.com/juanpoggi12/JGNSolutions/models"
)

func ConvertUserProfileRequestToModel(req dto.UserSearchRequest) models.UserProfile {
	return models.UserProfile{
		UserID:     ToObjectID(req.UserID),
		Name:       req.Name,
		BirthDate:  req.BirthDate,
		Weight:     req.Weight,
		Height:     req.Height,
		Experience: req.Experience,
		Goals:      req.Goals,
	}
}

func ConvertUserProfileModelToResponse(up models.UserProfile) dto.UserProfileResponse {
	return dto.UserProfileResponse{
		ID:         up.ID.Hex(),
		UserID:     up.UserID.Hex(),
		Name:       up.Name,
		BirthDate:  up.BirthDate,
		Weight:     up.Weight,
		Height:     up.Height,
		Experience: up.Experience,
		Goals:      up.Goals,
	}
}
