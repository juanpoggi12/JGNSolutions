package utils

import (
	"JGNSolutions/backend/dto"
	"JGNSolutions/backend/models"
)

func ConvertRoutineRequestToModel(req dto.RoutineRequest) models.Routine {
	return models.Routine{
		Name:        req.Name,
		Description: req.Description,
		UserID:      ToObjectID(req.UserID),
		Exercises:   []models.RoutineExercise{}, // se cargan aparte
	}
}

func ConvertRoutineModelToResponse(r models.Routine) dto.RoutineResponse {
	return dto.RoutineResponse{
		ID:          r.ID.Hex(),
		Name:        r.Name,
		Description: r.Description,
		UserID:      r.UserID.Hex(),
	}
}
