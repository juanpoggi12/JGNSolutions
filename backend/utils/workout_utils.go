package utils

import (
	"JGNSolutions/backend/dto"
	"JGNSolutions/backend/models"
)

func ConvertWorkoutRequestToModel(req dto.WorkoutRequest) models.Workout {
	return models.Workout{
		UserID:    ToObjectID(req.UserID),
		RoutineID: ToObjectID(req.RoutineID),
		Date:      req.Date,
		Duration:  req.Duration,
	}
}

func ConvertWorkoutModelToResponse(w models.Workout) dto.WorkoutResponse {
	return dto.WorkoutResponse{
		ID:        w.ID.Hex(),
		UserID:    w.UserID.Hex(),
		RoutineID: w.RoutineID.Hex(),
		Date:      w.Date,
		Duration:  w.Duration,
	}
}
