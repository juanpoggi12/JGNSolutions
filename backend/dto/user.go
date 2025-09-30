package dto

type UserCreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=40"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=72"`
	Role     string `json:"role" binding:"omitempty,oneof=ADMIN USER"` // default USER en servicio
	IsActive *bool  `json:"is_active" binding:"omitempty"`
}

type UserUpdateRequest struct {
	Username *string `json:"username" binding:"omitempty,min=3,max=40"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty,min=6,max=72"` // si permite reset
	Role     *string `json:"role" binding:"omitempty,oneof=ADMIN USER"`
	IsActive *bool   `json:"is_active" binding:"omitempty"`
}

type UserSearchRequest struct {
	Username string `form:"username"`
	Email    string `form:"email"`
	Role     string `form:"role" binding:"omitempty,oneof=ADMIN USER"`
	IsActive *bool  `form:"is_active"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"` // RFC3339
	UpdatedAt string `json:"updated_at"` // RFC3339
	IsActive  bool   `json:"is_active"`
}
