package utils

import (
	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"
)

// Convierte DTO de creación a modelo (Request -> Model)
func ConvertRoutineExerciseCreateRequestToModel(req dto.RoutineExerciseCreateRequest) models.RoutineExercise {
	return models.RoutineExercise{
		RoutineID:    ToObjectID(req.RoutineID),
		ExerciseID:   ToObjectID(req.ExerciseID),
		Orden:        req.Orden,
		Series:       req.Series,
		Repeticiones: req.Repeticiones,
		PesoObjetivo: req.PesoObjetivo,
		TiempoObjSeg: req.TiempoObjSeg,
		Notas:        req.Notas,
	}
}

// Convierte modelo a DTO de respuesta (Model -> Response)
func ConvertRoutineExerciseModelToResponse(re models.RoutineExercise) dto.RoutineExerciseResponse {
	return dto.RoutineExerciseResponse{
		ID:           re.ID.Hex(),
		RoutineID:    re.RoutineID.Hex(),
		ExerciseID:   re.ExerciseID.Hex(),
		Orden:        re.Orden,
		Series:       re.Series,
		Repeticiones: re.Repeticiones,
		PesoObjetivo: re.PesoObjetivo,
		TiempoObjSeg: re.TiempoObjSeg,
		Notas:        re.Notas,
	}
}
