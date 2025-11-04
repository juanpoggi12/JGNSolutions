package dto

type WorkoutEntryCreateRequest struct {
	WorkoutSessionID string   `json:"workout_session_id" binding:"required,hexadecimal,len=24"`
	ExerciseID       string   `json:"exercise_id" binding:"required,hexadecimal,len=24"`
	SetNumber        int      `json:"set_number" binding:"required,min=1"`
	RepsDone         *int     `json:"reps_done,omitempty" binding:"omitempty,min=0,max=1000"`
	WeightUsed       *float64 `json:"weight_used,omitempty" binding:"omitempty,gte=0"`
	TimeSec          *int     `json:"time_sec,omitempty" binding:"omitempty,gte=0"`
	PerceivedEffort  *int     `json:"perceived_effort,omitempty" binding:"omitempty,min=0,max=10"`
}

type WorkoutEntrySearchRequest struct {
	WorkoutSessionID string `form:"workout_session_id" binding:"omitempty,hexadecimal,len=24"`
	ExerciseID       string `form:"exercise_id" binding:"omitempty,hexadecimal,len=24"`
}

type WorkoutEntryResponse struct {
	ID               string   `json:"id"`
	WorkoutSessionID string   `json:"workout_session_id"`
	ExerciseID       string   `json:"exercise_id"`
	SetNumber        int      `json:"set_number"`
	RepsDone         *int     `json:"reps_done,omitempty"`
	WeightUsed       *float64 `json:"weight_used,omitempty"`
	TimeSec          *int     `json:"time_sec,omitempty"`
	PerceivedEffort  *int     `json:"perceived_effort,omitempty"`
}
