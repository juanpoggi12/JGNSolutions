package utils

import (
	"github.com/juanpoggi12/JGNSolutions/dto"
	"github.com/juanpoggi12/JGNSolutions/models"
)

func ConvertWorkoutSessionCreateRequestToModel(req dto.WorkoutSessionCreateRequest) models.WorkoutSession {
	// Note: DTO uses StartTime/EndTime as strings (RFC3339); mapper should parse to time.Time in real code
	return models.WorkoutSession{
		UserID:    ToObjectID(req.UserID),
		RoutineID: nil, // set outside if provided
	}
}

func ConvertWorkoutSessionModelToResponse(ws models.WorkoutSession) dto.WorkoutSessionResponse {
	return dto.WorkoutSessionResponse{
		ID:        ws.ID.Hex(),
		UserID:    ws.UserID.Hex(),
		RoutineID: nil,
		// StartTime/EndTime should be formatted in mapper
	}
}
