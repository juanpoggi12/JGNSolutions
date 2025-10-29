package utils

import (
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
)

// Convierte DTO de creación a modelo (Request -> Model)
// Ya no usa UserID del body: el userID se asigna externamente (desde el Actor)
func ConvertRoutineCreateRequestToModel(req dto.RoutineCreateRequest) (models.Routine, error) {
	return models.Routine{
		Name:        req.Name,
		Description: req.Description,
		IsTemplate:  req.IsTemplate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsDeleted:   false,
	}, nil
}

// Convierte modelo a DTO de respuesta (Model -> Response)
func ConvertRoutineModelToResponse(r models.Routine) dto.RoutineResponse {
	return dto.RoutineResponse{
		ID:          r.ID.Hex(),
		Name:        r.Name,
		Description: r.Description,
		IsTemplate:  r.IsTemplate,
		CreatedAt:   r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   r.UpdatedAt.Format(time.RFC3339),
		IsDeleted:   r.IsDeleted,
	}
}

// Convierte RoutineSearchRequest en filtro Mongo
// Ya no usa UserID del body; si se necesita filtrar por usuario, se pasa externamente.
func BuildRoutineSearchFilter(search dto.RoutineSearchRequest, userID string) bson.M {
	filter := bson.M{}

	if userID != "" {
		filter["userId"] = ToObjectID(userID)
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
