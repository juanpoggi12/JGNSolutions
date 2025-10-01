package utils

import (
	"github.com/juanpoggi12/JGNSolutions/dto"
	"github.com/juanpoggi12/JGNSolutions/models"
)

func ConvertRoutineCreateRequestToModel(req dto.RoutineCreateRequest) models.Routine {
	return models.Routine{
		Name:        req.Name,
		Description: req.Description,
		UserID:      ToObjectID(req.UserID),
		// Exercises loaded separately
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
