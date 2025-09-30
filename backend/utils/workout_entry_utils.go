package utils

import (
	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"
)

// Request -> Model
func ConvertWorkoutEntryCreateRequestToModel(req dto.WorkoutEntryCreateRequest) models.WorkoutEntry {
	return models.WorkoutEntry{
		WorkoutSessionID:   ToObjectID(req.WorkoutSessionID),
		ExerciseID:         ToObjectID(req.ExerciseID),
		Serie:              req.Serie,
		RepsHechas:         req.RepsHechas,
		PesoUsado:          req.PesoUsado,
		TiempoSeg:          req.TiempoSeg,
		PercepcionEsfuerzo: req.PercepcionEsfuerzo,
	}
}

// Model -> Response
func ConvertWorkoutEntryModelToResponse(we models.WorkoutEntry) dto.WorkoutEntryResponse {
	return dto.WorkoutEntryResponse{
		ID:                 we.ID.Hex(),
		WorkoutSessionID:   we.WorkoutSessionID.Hex(),
		ExerciseID:         we.ExerciseID.Hex(),
		Serie:              we.Serie,
		RepsHechas:         we.RepsHechas,
		PesoUsado:          we.PesoUsado,
		TiempoSeg:          we.TiempoSeg,
		PercepcionEsfuerzo: we.PercepcionEsfuerzo,
	}
}
