package dto

type WorkoutSessionCreateRequest struct {
	UserID          string  `json:"user_id" binding:"required,hexadecimal,len=24"`
	RoutineID       *string `json:"routine_id" binding:"omitempty,hexadecimal,len=24"`
	FechaHoraInicio string  `json:"fecha_hora_inicio" binding:"required,datetime=2006-01-02T15:04:05Z07:00"` // RFC3339
	FechaHoraFin    *string `json:"fecha_hora_fin" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	NotasGenerales  string  `json:"notas_generales" binding:"omitempty,max=2000"`
}

type WorkoutSessionUpdateRequest struct {
	RoutineID       *string `json:"routine_id" binding:"omitempty,hexadecimal,len=24"`
	FechaHoraInicio *string `json:"fecha_hora_inicio" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	FechaHoraFin    *string `json:"fecha_hora_fin" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	NotasGenerales  *string `json:"notas_generales" binding:"omitempty,max=2000"`
}

type WorkoutSessionSearchRequest struct {
	UserID        string `form:"user_id" binding:"omitempty,hexadecimal,len=24"`
	RoutineID     string `form:"routine_id" binding:"omitempty,hexadecimal,len=24"`
	StartedAfter  string `form:"started_after" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	StartedBefore string `form:"started_before" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type WorkoutSessionResponse struct {
	ID              string  `json:"id"`
	UserID          string  `json:"user_id"`
	RoutineID       *string `json:"routine_id,omitempty"`
	FechaHoraInicio string  `json:"fecha_hora_inicio"` // RFC3339
	FechaHoraFin    string  `json:"fecha_hora_fin"`    // RFC3339
	NotasGenerales  string  `json:"notas_generales,omitempty"`
	CreatedAt       string  `json:"created_at"` // RFC3339
}
