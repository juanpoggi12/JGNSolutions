package utils

import (
	"JGNSolutions/backend/dto"
	"JGNSolutions/backend/models"
	"strings"
)

//model -> response

func ConvertExcersiceToModelResponse(e models.Excercise) dto.ExerciseResponse {

	return dto.ExcerciseResponse{
		ID:          e.ID.Hex(),
		Name:        e.Name,
		Description: e.Description,
		Category:    e.Category,
		MuscleGroup: e.MuscleGroup,
		Difficulty:  e.Difficulty,
		MediaURL:    e.MediaURL,
		Steps:       e.Steps,
		CreatedBy:   e.CreatedBy.Hex(),
	}
}

// Request -> Model
func ConvertExerciseRequestToModel(req dto.ExerciseRequest) models.Exercise {
	return models.Exercise{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		MuscleGroup: req.MuscleGroup,
		Difficulty:  req.Difficulty,
		MediaURL:    req.MediaURL,
		Steps:       req.Steps,
	}
}

// Filtro de bÃºsqueda
func MatchesExerciseSearch(e models.Exercise, search dto.ExerciseSearchRequest) bool {
	if search.Name != "" && !strings.Contains(strings.ToLower(e.Name), strings.ToLower(search.Name)) {
		return false
	}
	if search.Category != "" && !strings.Contains(strings.ToLower(e.Category), strings.ToLower(search.Category)) {
		return false
	}
	if search.MuscleGroup != "" && !strings.Contains(strings.ToLower(e.MuscleGroup), strings.ToLower(search.MuscleGroup)) {
		return false
	}
	if search.Difficulty != "" && !strings.Contains(strings.ToLower(e.Difficulty), strings.ToLower(search.Difficulty)) {
		return false
	}
	return true
}
