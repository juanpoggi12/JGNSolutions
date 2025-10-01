package dto

type RoutineExerciseCreateRequest struct {
	RoutineID     string   `json:"routine_id" binding:"required,hexadecimal,len=24"`
	ExerciseID    string   `json:"exercise_id" binding:"required,hexadecimal,len=24"`
	Order         int      `json:"order" binding:"required,min=1"`
	Sets          int      `json:"sets" binding:"required,min=1,max=50"`
	Reps          int      `json:"reps" binding:"required,min=0,max=1000"`
	TargetWeight  *float64 `json:"target_weight,omitempty" binding:"omitempty,gte=0"`
	TargetTimeSec *int     `json:"target_time_sec,omitempty" binding:"omitempty,gte=0"`
	Notes         string   `json:"notes,omitempty" binding:"omitempty,max=1000"`
}

type RoutineExerciseUpdateRequest struct {
	Order         *int     `json:"order" binding:"omitempty,min=1"`
	Sets          *int     `json:"sets" binding:"omitempty,min=1,max=50"`
	Reps          *int     `json:"reps" binding:"omitempty,min=0,max=1000"`
	TargetWeight  *float64 `json:"target_weight,omitempty" binding:"omitempty,gte=0"`
	TargetTimeSec *int     `json:"target_time_sec,omitempty" binding:"omitempty,gte=0"`
	Notes         *string  `json:"notes,omitempty" binding:"omitempty,max=1000"`
	// routine_id and exercise_id are not changed in update
}

type RoutineExerciseSearchRequest struct {
	RoutineID  string `form:"routine_id" binding:"omitempty,hexadecimal,len=24"`
	ExerciseID string `form:"exercise_id" binding:"omitempty,hexadecimal,len=24"`
}

type RoutineExerciseResponse struct {
	ID            string   `json:"id"`
	RoutineID     string   `json:"routine_id"`
	ExerciseID    string   `json:"exercise_id"`
	Order         int      `json:"order"`
	Sets          int      `json:"sets"`
	Reps          int      `json:"reps"`
	TargetWeight  *float64 `json:"target_weight,omitempty"`
	TargetTimeSec *int     `json:"target_time_sec,omitempty"`
	Notes         string   `json:"notes,omitempty"`
}
