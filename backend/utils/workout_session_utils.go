package utils

import (
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateRequest → Model
// Ya no toma el UserID del DTO; se asigna externamente en el service con actor.UserID
func ConvertWorkoutSessionCreateRequestToModel(req dto.WorkoutSessionCreateRequest) (models.WorkoutSession, error) {
	var routineID *primitive.ObjectID
	if req.RoutineID != nil {
		id, err := primitive.ObjectIDFromHex(*req.RoutineID)
		if err != nil {
			return models.WorkoutSession{}, err
		}
		routineID = &id
	}

	start, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return models.WorkoutSession{}, err
	}

	var end time.Time
	if req.EndTime != nil {
		end, err = time.Parse(time.RFC3339, *req.EndTime)
		if err != nil {
			return models.WorkoutSession{}, err
		}
	}

	return models.WorkoutSession{
		// UserID se asignará en el service
		RoutineID: routineID,
		StartTime: start,
		EndTime:   end,
		Notes:     req.Notes,
		CreatedAt: time.Now(),
	}, nil
}

// UpdateRequest → aplica cambios sobre el modelo existente
func ApplyWorkoutSessionUpdateToModel(ws *models.WorkoutSession, req dto.WorkoutSessionUpdateRequest) error {
	if req.RoutineID != nil {
		id := ToObjectID(*req.RoutineID)
		ws.RoutineID = &id
	}
	if req.StartTime != nil {
		start, err := time.Parse(time.RFC3339, *req.StartTime)
		if err != nil {
			return err
		}
		ws.StartTime = start
	}
	if req.EndTime != nil {
		end, err := time.Parse(time.RFC3339, *req.EndTime)
		if err != nil {
			return err
		}
		ws.EndTime = end
	}
	if req.Notes != nil {
		ws.Notes = *req.Notes
	}
	return nil
}

// Model → Response
// Quita UserID si el DTO ya no lo tiene
func ConvertWorkoutSessionModelToResponse(ws models.WorkoutSession) dto.WorkoutSessionResponse {
	var routineID *string
	if ws.RoutineID != nil {
		id := ws.RoutineID.Hex()
		routineID = &id
	}

	return dto.WorkoutSessionResponse{
		ID:              ws.ID.Hex(),
		RoutineID:       routineID,
		FechaHoraInicio: ws.StartTime.Format(time.RFC3339),
		FechaHoraFin:    ws.EndTime.Format(time.RFC3339),
		NotasGenerales:  ws.Notes,
		CreatedAt:       ws.CreatedAt.Format(time.RFC3339),
	}
}

// SearchRequest → filtro Mongo
// Ya no usa UserID del body; el filtro por usuario lo manejará el service usando actor.UserID
func BuildWorkoutSessionSearchFilter(search dto.WorkoutSessionSearchRequest, userID string) (bson.M, error) {
	filter := bson.M{}

	if userID != "" {
		filter["userId"] = ToObjectID(userID)
	}
	if search.RoutineID != "" {
		filter["routineId"] = ToObjectID(search.RoutineID)
	}
	if search.StartedAfter != "" {
		after, err := time.Parse(time.RFC3339, search.StartedAfter)
		if err != nil {
			return nil, err
		}
		filter["startTime"] = bson.M{"$gte": after}
	}
	if search.StartedBefore != "" {
		before, err := time.Parse(time.RFC3339, search.StartedBefore)
		if err != nil {
			return nil, err
		}
		if f, ok := filter["startTime"].(bson.M); ok {
			f["$lte"] = before
		} else {
			filter["startTime"] = bson.M{"$lte": before}
		}
	}
	return filter, nil
}
