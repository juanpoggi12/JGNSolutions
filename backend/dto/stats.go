package dto

type ExerciseStatResponse struct {
	ExerciseID string `json:"exercise_id"`
	Name       string `json:"name"`
	UsageCount int    `json:"usage_count"`
}

type RoutineStatResponse struct {
	RoutineID  string `json:"routine_id"`
	Name       string `json:"name"`
	UsageCount int    `json:"usage_count"`
}
