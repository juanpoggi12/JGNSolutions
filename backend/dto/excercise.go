package dto

type ExerciseCreateRequest struct {
	Nombre          string   `json:"nombre" binding:"required,min=2,max=120"`
	Descripcion     string   `json:"descripcion" binding:"omitempty,max=2000"`
	Categoria       string   `json:"categoria" binding:"required,oneof=FUERZA CARDIO FLEXIBILIDAD OTRA"`
	GrupoMuscular   string   `json:"grupo_muscular" binding:"required,oneof=PECHO ESPALDA PIERNA HOMBRO BRAZO CORE"`
	Dificultad      string   `json:"dificultad" binding:"required,oneof=BAJA MEDIA ALTA"`
	MediaURL        string   `json:"media_url" binding:"omitempty,url"`
	Instrucciones   []string `json:"instrucciones" binding:"omitempty,max=100"`
	CreatedByUserID string   `json:"created_by_user_id" binding:"required,hexadecimal,len=24"`
}

type ExerciseUpdateRequest struct {
	// PATCH: usar punteros para distinguir "no enviado" de "valor cero"
	Nombre        *string   `json:"nombre" binding:"omitempty,min=2,max=120"`
	Descripcion   *string   `json:"descripcion" binding:"omitempty,max=2000"`
	Categoria     *string   `json:"categoria" binding:"omitempty,oneof=FUERZA CARDIO FLEXIBILIDAD OTRA"`
	GrupoMuscular *string   `json:"grupo_muscular" binding:"omitempty,oneof=PECHO ESPALDA PIERNA HOMBRO BRAZO CORE"`
	Dificultad    *string   `json:"dificultad" binding:"omitempty,oneof=BAJA MEDIA ALTA"`
	MediaURL      *string   `json:"media_url" binding:"omitempty,url"`
	Instrucciones *[]string `json:"instrucciones" binding:"omitempty,max=100"`
	// created_by_user_id no se deberÃ­a cambiar normalmente
}

type ExerciseSearchRequest struct {
	Nombre        string `form:"nombre"`
	Categoria     string `form:"categoria" binding:"omitempty,oneof=FUERZA CARDIO FLEXIBILIDAD OTRA"`
	GrupoMuscular string `form:"grupo_muscular" binding:"omitempty,oneof=PECHO ESPALDA PIERNA HOMBRO BRAZO CORE"`
	Dificultad    string `form:"dificultad" binding:"omitempty,oneof=BAJA MEDIA ALTA"`
	CreatedBy     string `form:"created_by_user_id" binding:"omitempty,hexadecimal,len=24"`
	IncludeDel    bool   `form:"include_deleted"`
}

type ExerciseResponse struct {
	ID              string   `json:"id"`
	Nombre          string   `json:"nombre"`
	Descripcion     string   `json:"descripcion,omitempty"`
	Categoria       string   `json:"categoria"`
	GrupoMuscular   string   `json:"grupo_muscular"`
	Dificultad      string   `json:"dificultad"`
	MediaURL        string   `json:"media_url,omitempty"`
	Instrucciones   []string `json:"instrucciones,omitempty"`
	CreatedByUserID string   `json:"created_by_user_id"`
	CreatedAt       string   `json:"created_at"` // RFC3339 en el mapper
	UpdatedAt       string   `json:"updated_at"`
	IsDeleted       bool     `json:"is_deleted"`
}
