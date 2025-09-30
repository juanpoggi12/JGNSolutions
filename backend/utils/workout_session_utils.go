package utils

import (
	"JGNSolutions/backend/dto"
	"JGNSolutions/backend/models"
)

func ConvertWorkoutSessionRequestToModel(req dto.WorkoutSessionRequest) models.WorkoutSession {
	return models.WorkoutSession{
		UserID:  ToObjectID(req.UserID),
		Date:    req.Date,
		Entries: []models.WorkoutEntry{}, // se cargan aparte
	}
}

func ConvertWorkoutSessionModelToResponse(ws models.WorkoutSession) dto.WorkoutSessionResponse {
	return dto.WorkoutSessionResponse{
		ID:     ws.ID.Hex(),
		UserID: ws.UserID.Hex(),
		Date:   ws.Date,
	}
}
