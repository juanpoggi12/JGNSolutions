package utils

import (
	"time"

    "github.com/juanpoggi12/JGNSolutions/backend/dto"
    "github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

// Convierte DTO de creación a modelo (Request -> Model)
func ConvertRoutineCreateRequestToModel(req dto.RoutineCreateRequest) (models.Routine, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return models.Routine{}, err
	}

	return models.Routine{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		IsTemplate:  req.IsTemplate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsDeleted:   false,
	}, nil
}

// Aplica cambios de un UpdateRequest sobre un modelo existente (PATCH)
func ApplyRoutineUpdateToModel(r *models.Routine, req dto.RoutineUpdateRequest) {
	if req.Name != nil {
		r.Name = *req.Name
	}
	if req.Description != nil {
		r.Description = *req.Description
	}
	if req.IsTemplate != nil {
		r.IsTemplate = *req.IsTemplate
	}
	r.UpdatedAt = time.Now()
}

// Convierte modelo a DTO de respuesta (Model -> Response)
func ConvertRoutineModelToResponse(r models.Routine) dto.RoutineResponse {
	return dto.RoutineResponse{
		ID:          r.ID.Hex(),
		UserID:      r.UserID.Hex(),
		Name:        r.Name,
		Description: r.Description,
		IsTemplate:  r.IsTemplate,
		CreatedAt:   r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   r.UpdatedAt.Format(time.RFC3339),
		IsDeleted:   r.IsDeleted,
	}
}

// Convierte RoutineSearchRequest en filtro Mongo
func BuildRoutineSearchFilter(search dto.RoutineSearchRequest) bson.M {
	filter := bson.M{}
	if search.UserID != "" {
		filter["userId"] = ToObjectID(search.UserID)
	}
	if search.Name != "" {
		filter["name"] = bson.M{"$regex": search.Name, "$options": "i"}
	}
	if search.IsTemplate != nil {
		filter["isTemplate"] = *search.IsTemplate
	}
	if !search.IncludeDel {
		filter["isDeleted"] = false
	}
	return filter
}
