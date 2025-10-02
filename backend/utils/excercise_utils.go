package utils

import (
	"time"

    "github.com/juanpoggi12/JGNSolutions/backend/dto"
    "github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ----------------- Model → DTO -----------------

func ConvertExerciseModelToResponse(e models.Exercise) dto.ExerciseResponse {
	return dto.ExerciseResponse{
		ID:              e.ID.Hex(),
		Name:            e.Name,
		Description:     e.Description,
		Category:        string(e.Category),
		MuscleGroup:     string(e.MuscleGroup),
		Difficulty:      string(e.Difficulty),
		MediaURL:        e.MediaURL,
		Instructions:    e.Instructions,
		CreatedByUserID: e.CreatedByUserID.Hex(),
		CreatedAt:       e.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       e.UpdatedAt.Format(time.RFC3339),
		IsDeleted:       e.IsDeleted,
	}
}

// ----------------- DTO → Model -----------------

func ConvertExerciseCreateRequestToModel(req dto.ExerciseCreateRequest) (models.Exercise, error) {
	createdBy, err := primitive.ObjectIDFromHex(req.CreatedByUserID)
	if err != nil {
		return models.Exercise{}, err
	}
	return models.Exercise{
		Name:            req.Name,
		Description:     req.Description,
		Category:        models.ExerciseCategory(req.Category),
		MuscleGroup:     models.MuscleGroup(req.MuscleGroup),
		Difficulty:      models.Difficulty(req.Difficulty),
		MediaURL:        req.MediaURL,
		Instructions:    req.Instructions,
		CreatedByUserID: createdBy,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		IsDeleted:       false,
	}, nil
}

func ApplyExerciseUpdateToModel(e *models.Exercise, req dto.ExerciseUpdateRequest) {
	if req.Name != nil {
		e.Name = *req.Name
	}
	if req.Description != nil {
		e.Description = *req.Description
	}
	if req.Category != nil {
		e.Category = models.ExerciseCategory(*req.Category)
	}
	if req.MuscleGroup != nil {
		e.MuscleGroup = models.MuscleGroup(*req.MuscleGroup)
	}
	if req.Difficulty != nil {
		e.Difficulty = models.Difficulty(*req.Difficulty)
	}
	if req.MediaURL != nil {
		e.MediaURL = *req.MediaURL
	}
	if req.Instructions != nil {
		e.Instructions = *req.Instructions
	}
	e.UpdatedAt = time.Now()
}

// ----------------- SearchRequest → Filtro Mongo -----------------

func BuildExerciseSearchFilter(search dto.ExerciseSearchRequest) bson.M {
	filter := bson.M{}
	if search.Name != "" {
		filter["name"] = bson.M{"$regex": search.Name, "$options": "i"}
	}
	if search.Category != "" {
		filter["category"] = search.Category
	}
	if search.MuscleGroup != "" {
		filter["muscleGroup"] = search.MuscleGroup
	}
	if search.Difficulty != "" {
		filter["difficulty"] = search.Difficulty
	}
	if search.CreatedBy != "" {
		filter["createdByUserId"] = ToObjectID(search.CreatedBy)
	}
	if !search.IncludeDel {
		filter["isDeleted"] = false
	}
	return filter
}
