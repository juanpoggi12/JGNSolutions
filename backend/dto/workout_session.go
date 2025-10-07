package dto

type WorkoutSessionCreateRequest struct {
	RoutineID *string `json:"routine_id" binding:"omitempty,hexadecimal,len=24"`
	StartTime string  `json:"start_time" binding:"required,datetime=2006-01-02T15:04:05Z07:00"` // RFC3339
	EndTime   *string `json:"end_time" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Notes     string  `json:"notes" binding:"omitempty,max=2000"`
}

type WorkoutSessionUpdateRequest struct {
	RoutineID *string `json:"routine_id" binding:"omitempty,hexadecimal,len=24"`
	StartTime *string `json:"start_time" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	EndTime   *string `json:"end_time" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Notes     *string `json:"notes" binding:"omitempty,max=2000"`
}

type WorkoutSessionSearchRequest struct {
	RoutineID     string `form:"routine_id" binding:"omitempty,hexadecimal,len=24"`
	StartedAfter  string `form:"started_after" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	StartedBefore string `form:"started_before" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type WorkoutSessionResponse struct {
	ID              string  `json:"id"`
	RoutineID       *string `json:"routine_id,omitempty"`
	FechaHoraInicio string  `json:"fecha_hora_inicio"` // RFC3339
	FechaHoraFin    string  `json:"fecha_hora_fin"`    // RFC3339
	NotasGenerales  string  `json:"notas_generales,omitempty"`
	CreatedAt       string  `json:"created_at"` // RFC3339
}
