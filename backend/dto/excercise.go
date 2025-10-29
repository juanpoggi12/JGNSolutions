package dto

type ExerciseCreateRequest struct {
	Name            string   `json:"name" binding:"required,min=2,max=120"`
	Description     string   `json:"description" binding:"omitempty,max=2000"`
	Category        string   `json:"category" binding:"required,oneof=FUERZA CARDIO FLEXIBILIDAD OTRA"`
	MuscleGroup     string   `json:"muscle_group" binding:"required,oneof=PECHO ESPALDA PIERNA HOMBRO BRAZO CORE"`
	Difficulty      string   `json:"difficulty" binding:"required,oneof=BAJA MEDIA ALTA"`
	MediaURL        string   `json:"media_url" binding:"omitempty,url"`
	Instructions    []string `json:"instructions" binding:"omitempty,max=100"`
	CreatedByUserID string   `json:"created_by_user_id" binding:"required,hexadecimal,len=24"`
}

type ExerciseUpdateRequest struct {
	Name         *string   `json:"name" binding:"omitempty,min=2,max=120"`
	Description  *string   `json:"description" binding:"omitempty,max=2000"`
	Category     *string   `json:"category" binding:"omitempty,oneof=FUERZA CARDIO FLEXIBILIDAD OTRA"`
	MuscleGroup  *string   `json:"muscle_group" binding:"omitempty,oneof=PECHO ESPALDA PIERNA HOMBRO BRAZO CORE"`
	Difficulty   *string   `json:"difficulty" binding:"omitempty,oneof=BAJA MEDIA ALTA"`
	MediaURL     *string   `json:"media_url" binding:"omitempty,url"`
	Instructions *[]string `json:"instructions" binding:"omitempty,max=100"`
}

type ExerciseResponse struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description,omitempty"`
	Category     string   `json:"category"`
	MuscleGroup  string   `json:"muscle_group"`
	Difficulty   string   `json:"difficulty"`
	MediaURL     string   `json:"media_url,omitempty"`
	Instructions []string `json:"instructions,omitempty"`
	CreatedBy    string   `json:"created_by"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
	IsDeleted    bool     `json:"is_deleted"`
}

type ExerciseCatalogQuery struct {
	Q           string `form:"q"`
	Category    string `form:"category"`
	MuscleGroup string `form:"muscleGroup"`
	Difficulty  string `form:"difficulty"`
	Page        int    `form:"page"`
	Limit       int    `form:"limit"`
}

type ExerciseCatalogResponse struct {
	Items []ExerciseResponse `json:"items"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
}
