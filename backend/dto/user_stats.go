package dto

type WorkoutSummaryResponse struct {
	ByExercise         []WorkoutSummaryExercise `json:"byExercise"`
	ByPeriod           []WorkoutSummaryPeriod   `json:"byPeriod"`
	PRs                []WorkoutPersonalRecord  `json:"prs"`
	ImprovementPercent *float64                 `json:"improvementPercent,omitempty"` // <-- AÑADIR ESTA LÍNEA
}

type WorkoutSummaryExercise struct {
	ExerciseID   string  `json:"exerciseId"`
	Name         string  `json:"name"`
	MuscleGroup  string  `json:"muscleGroup"`
	TotalSets    int     `json:"totalSets"`
	TotalReps    int     `json:"totalReps"`
	TotalWeight  float64 `json:"totalWeight"`
	TotalTimeSec int     `json:"totalTimeSec"`
}

type WorkoutSummaryPeriod struct {
	Period   string `json:"period"`
	Label    string `json:"label"`
	Sessions int    `json:"sessions"`
}

type WorkoutPersonalRecord struct {
	ExerciseID string  `json:"exerciseId"`
	Name       string  `json:"name"`
	MaxWeight  float64 `json:"maxWeight"`
	MaxReps    int     `json:"maxReps"`
	MaxTimeSec int     `json:"maxTimeSec"`
}

// --- User frequency ---

type UserFrequencyResponse struct {
	Buckets []UserFrequencyBucket `json:"buckets"`
	Total   int                   `json:"total"`
}

type UserFrequencyBucket struct {
	Period string `json:"period"`
	Label  string `json:"label"`
	Count  int    `json:"count"`
}

// --- Top routines ---

type UserTopRoutine struct {
	RoutineID string `json:"routineId"`
	Name      string `json:"name"`
	Uses      int    `json:"uses"`
}

// --- Exercise progress ---

type UserExerciseProgressResponse struct {
	Labels  []string  `json:"labels"`
	Sets    []int     `json:"sets"`
	Reps    []int     `json:"reps"`
	Weight  []float64 `json:"weight"`
	TimeSec []int     `json:"timeSec"`
}
