package utils

import (
	"JGNSolutions/backend/dto"
	"JGNSolutions/backend/models"
)

func ConvertWorkoutEntryRequestToModel(req dto.WorkoutEntryRequest) models.WorkoutEntry {
	return models.WorkoutEntry{
		WorkoutID:  ToObjectID(req.WorkoutID),
		ExerciseID: ToObjectID(req.ExerciseID),
		Sets:       req.Sets,
		Reps:       req.Reps,
		Weight:     req.Weight,
	}
}

func ConvertWorkoutEntryModelToResponse(we models.WorkoutEntry) dto.WorkoutEntryResponse {
	return dto.WorkoutEntryResponse{
		ID:         we.ID.Hex(),
		WorkoutID:  we.WorkoutID.Hex(),
		ExerciseID: we.ExerciseID.Hex(),
		Sets:       we.Sets,
		Reps:       we.Reps,
		Weight:     we.Weight,
	}
}
