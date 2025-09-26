package dto

type WorkoutEntryCreateRequest struct {
	WorkoutSessionID   string   `json:"workout_session_id" binding:"required,hexadecimal,len=24"`
	ExerciseID         string   `json:"exercise_id" binding:"required,hexadecimal,len=24"`
	Serie              int      `json:"serie" binding:"required,min=1"`
	RepsHechas         *int     `json:"reps_hechas,omitempty" binding:"omitempty,min=0,max=1000"`
	PesoUsado          *float64 `json:"peso_usado,omitempty" binding:"omitempty,gte=0"`
	TiempoSeg          *int     `json:"tiempo_seg,omitempty" binding:"omitempty,gte=0"`
	PercepcionEsfuerzo *int     `json:"percepcion_esfuerzo,omitempty" binding:"omitempty,min=0,max=10"` // si usan RPE 0-10
}

type WorkoutEntryUpdateRequest struct {
	Serie              *int     `json:"serie" binding:"omitempty,min=1"`
	RepsHechas         *int     `json:"reps_hechas,omitempty" binding:"omitempty,min=0,max=1000"`
	PesoUsado          *float64 `json:"peso_usado,omitempty" binding:"omitempty,gte=0"`
	TiempoSeg          *int     `json:"tiempo_seg,omitempty" binding:"omitempty,gte=0"`
	PercepcionEsfuerzo *int     `json:"percepcion_esfuerzo,omitempty" binding:"omitempty,min=0,max=10"`
}

type WorkoutEntrySearchRequest struct {
	WorkoutSessionID string `form:"workout_session_id" binding:"omitempty,hexadecimal,len=24"`
	ExerciseID       string `form:"exercise_id" binding:"omitempty,hexadecimal,len=24"`
	PageQuery
}

type WorkoutEntryResponse struct {
	ID                 string   `json:"id"`
	WorkoutSessionID   string   `json:"workout_session_id"`
	ExerciseID         string   `json:"exercise_id"`
	Serie              int      `json:"serie"`
	RepsHechas         *int     `json:"reps_hechas,omitempty"`
	PesoUsado          *float64 `json:"peso_usado,omitempty"`
	TiempoSeg          *int     `json:"tiempo_seg,omitempty"`
	PercepcionEsfuerzo *int     `json:"percepcion_esfuerzo,omitempty"`
}
