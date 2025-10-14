package dto

// 📊 Estadísticas de ejercicios
type ExerciseStatResponse struct {
	ExerciseID string `json:"exercise_id"`
	Name       string `json:"name"`
	UsageCount int    `json:"usage_count"`
}

// 📊 Estadísticas de rutinas
type RoutineStatResponse struct {
	RoutineID  string `json:"routine_id"`
	Name       string `json:"name"`
	UsageCount int    `json:"usage_count"`
}

// 📊 Estadísticas de perfiles de usuario por nivel
type ProfileLevelStat struct {
	Level string `json:"level"`
	Count int    `json:"count"`
}

// 📊 Estadísticas de perfiles de usuario por objetivo
type ProfileGoalStat struct {
	Goal  string `json:"goal"`
	Count int    `json:"count"`
}
