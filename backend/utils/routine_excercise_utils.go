package utils

import (
    "github.com/juanpoggi12/JGNSolutions/backend/dto"
    "github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Convierte DTO de creación a modelo (Request -> Model)
func ConvertRoutineExerciseCreateRequestToModel(req dto.RoutineExerciseCreateRequest) (models.RoutineExercise, error) {
	routineID, err := primitive.ObjectIDFromHex(req.RoutineID)
	if err != nil {
		return models.RoutineExercise{}, err
	}
	exerciseID, err := primitive.ObjectIDFromHex(req.ExerciseID)
	if err != nil {
		return models.RoutineExercise{}, err
	}

	return models.RoutineExercise{
		RoutineID:     routineID,
		ExerciseID:    exerciseID,
		Order:         req.Order,
		Sets:          req.Sets,
		Reps:          req.Reps,
		TargetWeight:  req.TargetWeight,
		TargetTimeSec: req.TargetTimeSec,
		Notes:         req.Notes,
	}, nil
}

func ApplyRoutineExerciseUpdateToModel(re *models.RoutineExercise, req dto.RoutineExerciseUpdateRequest) {
	if req.Order != nil {
		re.Order = *req.Order
	}
	if req.Sets != nil {
		re.Sets = *req.Sets
	}
	if req.Reps != nil {
		re.Reps = *req.Reps
	}
	if req.TargetWeight != nil {
		re.TargetWeight = req.TargetWeight
	}
	if req.TargetTimeSec != nil {
		re.TargetTimeSec = req.TargetTimeSec
	}
	if req.Notes != nil {
		re.Notes = *req.Notes
	}
}

func ConvertRoutineExerciseModelToResponse(re models.RoutineExercise) dto.RoutineExerciseResponse {
	return dto.RoutineExerciseResponse{
		ID:            re.ID.Hex(),
		RoutineID:     re.RoutineID.Hex(),
		ExerciseID:    re.ExerciseID.Hex(),
		Order:         re.Order,
		Sets:          re.Sets,
		Reps:          re.Reps,
		TargetWeight:  re.TargetWeight,
		TargetTimeSec: re.TargetTimeSec,
		Notes:         re.Notes,
	}
}

func BuildRoutineExerciseSearchFilter(search dto.RoutineExerciseSearchRequest) bson.M {
	filter := bson.M{}
	if search.RoutineID != "" {
		filter["routineId"] = ToObjectID(search.RoutineID)
	}
	if search.ExerciseID != "" {
		filter["exerciseId"] = ToObjectID(search.ExerciseID)
	}
	return filter
}
