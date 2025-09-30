package utils

import (
	"JGNSolutions/backend/dto"
	"JGNSolutions/backend/models"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"
)

func ConvertRoutineExerciseRequestToModel(req dto.RoutineExerciseRequest) models.RoutineExercise {
	return models.RoutineExercise{
		ExerciseID: ToObjectID(req.ExerciseID),
		Sets:       req.Sets,
		Reps:       req.Reps,
		Weight:     req.Weight,
	}
}

func ConvertRoutineExerciseModelToResponse(re models.RoutineExercise) dto.RoutineExerciseResponse {
	return dto.RoutineExerciseResponse{
		ExerciseID: re.ExerciseID.Hex(),
		Sets:       re.Sets,
		Reps:       re.Reps,
		Weight:     re.Weight,
	}
}
