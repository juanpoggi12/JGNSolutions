package utils

import (
	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateRequest → Model
func ConvertWorkoutEntryCreateRequestToModel(req dto.WorkoutEntryCreateRequest) (models.WorkoutEntry, error) {
	sessionID, err := primitive.ObjectIDFromHex(req.WorkoutSessionID)
	if err != nil {
		return models.WorkoutEntry{}, err
	}
	exerciseID, err := primitive.ObjectIDFromHex(req.ExerciseID)
	if err != nil {
		return models.WorkoutEntry{}, err
	}

	return models.WorkoutEntry{
		WorkoutSessionID: sessionID,
		ExerciseID:       exerciseID,
		SetNumber:        req.SetNumber,
		RepsDone:         req.RepsDone,
		WeightUsed:       req.WeightUsed,
		TimeSec:          req.TimeSec,
		PerceivedEffort:  req.PerceivedEffort,
	}, nil
}

// Model → Response
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

// SearchRequest → filtro Mongo
func BuildWorkoutEntrySearchFilter(search dto.WorkoutEntrySearchRequest) bson.M {
	filter := bson.M{}
	if search.WorkoutSessionID != "" {
		filter["workoutSessionId"] = ToObjectID(search.WorkoutSessionID)
	}
	if search.ExerciseID != "" {
		filter["exerciseId"] = ToObjectID(search.ExerciseID)
	}
	return filter
}
