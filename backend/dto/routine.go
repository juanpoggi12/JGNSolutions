package dto


type RoutineCreateRequest struct {
	UserID      string `json:"user_id" binding:"required,hexadecimal,len=24"`
	Nombre      string `json:"nombre" binding:"required,min=2,max=120"`
	Descripcion string `json:"descripcion" binding:"omitempty,max=1000"`
	IsTemplate  bool   `json:"is_template"`
}

type RoutineUpdateRequest struct {
	Nombre      *string `json:"nombre" binding:"omitempty,min=2,max=120"`
	Descripcion *string `json:"descripcion" binding:"omitempty,max=1000"`
	IsTemplate  *bool   `json:"is_template" binding:"omitempty"`
}

type RoutineSearchRequest struct {
	UserID     string `form:"user_id" binding:"omitempty,hexadecimal,len=24"`
	Nombre     string `form:"nombre"`
	IsTemplate *bool  `form:"is_template"`
	IncludeDel bool   `form:"include_deleted"`
	PageQuery
}

type RoutineResponse struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion,omitempty"`
	IsTemplate  bool   `json:"is_template"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	IsDeleted   bool   `json:"is_deleted"`
}
