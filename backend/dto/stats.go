package dto

// Estadísticas de ejercicios
type ExerciseStatResponse struct {
	ExerciseID string `json:"exercise_id"`
	Name       string `json:"name"`
	UsageCount int    `json:"usage_count"`
}

// Estadísticas de rutinas
type RoutineStatResponse struct {
	RoutineID  string `json:"routine_id"`
	Name       string `json:"name"`
	UsageCount int    `json:"usage_count"`
}
