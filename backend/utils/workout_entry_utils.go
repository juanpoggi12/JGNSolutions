package utils

import (
	"github.com/juanpoggi12/JGNSolutions/dto"
	"github.com/juanpoggi12/JGNSolutions/models"
)

// Request -> Model
func ConvertWorkoutEntryCreateRequestToModel(req dto.WorkoutEntryCreateRequest) models.WorkoutEntry {
	return models.WorkoutEntry{
		WorkoutSessionID: ToObjectID(req.WorkoutSessionID),
		ExerciseID:       ToObjectID(req.ExerciseID),
		SetNumber:        req.SetNumber,
		RepsDone:         req.RepsDone,
		WeightUsed:       req.WeightUsed,
		TimeSec:          req.TimeSec,
		PerceivedEffort:  req.PerceivedEffort,
	}
}

// Model -> Response
func ConvertWorkoutEntryModelToResponse(we models.WorkoutEntry) dto.WorkoutEntryResponse {
	return dto.WorkoutEntryResponse{
		ID:               we.ID.Hex(),
		WorkoutSessionID: we.WorkoutSessionID.Hex(),
		ExerciseID:       we.ExerciseID.Hex(),
		SetNumber:        we.SetNumber,
		RepsDone:         we.RepsDone,
		WeightUsed:       we.WeightUsed,
		TimeSec:          we.TimeSec,
		PerceivedEffort:  we.PerceivedEffort,
	}
}
