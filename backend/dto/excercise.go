package dto

import (
	"time"

	"JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExerciseCreateRequest struct {
	Name            string                  `json:"name" binding:"required,min=2,max=120"`
	Description     string                  `json:"description" binding:"omitempty,max=2000"`
	Category        models.ExerciseCategory `json:"category" binding:"required,oneof=FUERZA CARDIO FLEXIBILIDAD OTRA"`
	MuscleGroup     models.MuscleGroup      `json:"muscle_group" binding:"required,oneof=PECHO ESPALDA PIERNA HOMBRO BRAZO CORE"`
	Difficulty      models.Difficulty       `json:"difficulty" binding:"required,oneof=BAJA MEDIA ALTA"`
	MediaURL        string                  `json:"media_url" binding:"omitempty,url"`
	Instructions    []string                `json:"instructions" binding:"omitempty,max=100"`
	CreatedByUserID primitive.ObjectID      `json:"created_by_user_id" binding:"required"`
}

type ExerciseUpdateRequest struct {
	Name         *string                  `json:"name" binding:"omitempty,min=2,max=120"`
	Description  *string                  `json:"description" binding:"omitempty,max=2000"`
	Category     *models.ExerciseCategory `json:"category" binding:"omitempty,oneof=FUERZA CARDIO FLEXIBILIDAD OTRA"`
	MuscleGroup  *models.MuscleGroup      `json:"muscle_group" binding:"omitempty,oneof=PECHO ESPALDA PIERNA HOMBRO BRAZO CORE"`
	Difficulty   *models.Difficulty       `json:"difficulty" binding:"omitempty,oneof=BAJA MEDIA ALTA"`
	MediaURL     *string                  `json:"media_url" binding:"omitempty,url"`
	Instructions *[]string                `json:"instructions" binding:"omitempty,max=100"`
}

type ExerciseSearchRequest struct {
	Name        string                  `form:"name"`
	Category    models.ExerciseCategory `form:"category" binding:"omitempty,oneof=FUERZA CARDIO FLEXIBILIDAD OTRA"`
	MuscleGroup models.MuscleGroup      `form:"muscle_group" binding:"omitempty,oneof=PECHO ESPALDA PIERNA HOMBRO BRAZO CORE"`
	Difficulty  models.Difficulty       `form:"difficulty" binding:"omitempty,oneof=BAJA MEDIA ALTA"`
	CreatedBy   primitive.ObjectID      `form:"created_by_user_id" binding:"omitempty"`
	IncludeDel  bool                    `form:"include_deleted"`
}

type ExerciseResponse struct {
	ID              primitive.ObjectID      `json:"id"`
	Name            string                  `json:"name"`
	Description     string                  `json:"description,omitempty"`
	Category        models.ExerciseCategory `json:"category"`
	MuscleGroup     models.MuscleGroup      `json:"muscle_group"`
	Difficulty      models.Difficulty       `json:"difficulty"`
	MediaURL        string                  `json:"media_url,omitempty"`
	Instructions    []string                `json:"instructions,omitempty"`
	CreatedByUserID primitive.ObjectID      `json:"created_by_user_id"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
	IsDeleted       bool                    `json:"is_deleted"`
}
