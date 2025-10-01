package utils

import (
	"github.com/juanpoggi12/JGNSolutions/dto"
	"github.com/juanpoggi12/JGNSolutions/models"
)

// Convierte DTO de creación a modelo (Request -> Model)
func ConvertRoutineExerciseCreateRequestToModel(req dto.RoutineExerciseCreateRequest) models.RoutineExercise {
	return models.RoutineExercise{
		RoutineID:     ToObjectID(req.RoutineID),
		ExerciseID:    ToObjectID(req.ExerciseID),
		Order:         req.Order,
		Sets:          req.Sets,
		Reps:          req.Reps,
		TargetWeight:  req.TargetWeight,
		TargetTimeSec: req.TargetTimeSec,
		Notes:         req.Notes,
	}
}

// Convierte modelo a DTO de respuesta (Model -> Response)
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
