package dto

type RoutineExerciseCreateRequest struct {
	RoutineID      string   `json:"routine_id" binding:"required,hexadecimal,len=24"`
	ExerciseID     string   `json:"exercise_id" binding:"required,hexadecimal,len=24"`
	Orden          int      `json:"orden" binding:"required,min=1"`
	Series         int      `json:"series" binding:"required,min=1,max=50"`
	Repeticiones   int      `json:"repeticiones" binding:"required,min=0,max=1000"`
	PesoObjetivo   *float64 `json:"peso_objetivo,omitempty" binding:"omitempty,gte=0"`
	TiempoObjSeg   *int     `json:"tiempo_objetivo_seg,omitempty" binding:"omitempty,gte=0"`
	Notas          string   `json:"notas,omitempty" binding:"omitempty,max=1000"`
}

type RoutineExerciseUpdateRequest struct {
	Orden        *int     `json:"orden" binding:"omitempty,min=1"`
	Series       *int     `json:"series" binding:"omitempty,min=1,max=50"`
	Repeticiones *int     `json:"repeticiones" binding:"omitempty,min=0,max=1000"`
	PesoObjetivo *float64 `json:"peso_objetivo,omitempty" binding:"omitempty,gte=0"`
	TiempoObjSeg *int     `json:"tiempo_objetivo_seg,omitempty" binding:"omitempty,gte=0"`
	Notas        *string  `json:"notas,omitempty" binding:"omitempty,max=1000"`
	// routine_id y exercise_id no se tocan en update
}

type RoutineExerciseSearchRequest struct {
	RoutineID string `form:"routine_id" binding:"omitempty,hexadecimal,len=24"`
	ExerciseID string `form:"exercise_id" binding:"omitempty,hexadecimal,len=24"`
	PageQuery
}

type RoutineExerciseResponse struct {
	ID            string   `json:"id"`
	RoutineID     string   `json:"routine_id"`
	ExerciseID    string   `json:"exercise_id"`
	Orden         int      `json:"orden"`
	Series        int      `json:"series"`
	Repeticiones  int      `json:"repeticiones"`
	PesoObjetivo  *float64 `json:"peso_objetivo,omitempty"`
	TiempoObjSeg  *int     `json:"tiempo_objetivo_seg,omitempty"`
	Notas         string   `json:"notas,omitempty"`
}