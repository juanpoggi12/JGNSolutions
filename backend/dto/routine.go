package dto

type RoutineCreateRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=120"`
	Description string `json:"description" binding:"omitempty,max=1000"`
	IsTemplate  bool   `json:"is_template"`
}

type RoutineUpdateRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=120"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
	IsTemplate  *bool   `json:"is_template" binding:"omitempty"`
}

type RoutineSearchRequest struct {
	Name       string `form:"name"`
	IsTemplate *bool  `form:"is_template"`
	IncludeDel bool   `form:"include_deleted"`
}

type RoutineResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsTemplate  bool   `json:"is_template"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	IsDeleted   bool   `json:"is_deleted"`
}
