package dto

// Crear una nueva sesión (cuando el usuario hace login con refresh token)
type SessionCreateRequest struct {
	UserID    string `json:"user_id" binding:"required,hexadecimal,len=24"`
	Token     string `json:"token" binding:"required,min=20"` // refresh token sin hashear
	ExpiresAt string `json:"expires_at" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

type SessionResponse struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	ExpiresAt string  `json:"expires_at"`           // RFC3339
	CreatedAt string  `json:"created_at"`           // RFC3339
	RevokedAt *string `json:"revoked_at,omitempty"` // RFC3339
}

type SessionIDUri struct {
	ID string `uri:"id" binding:"required,hexadecimal,len=24"`
}

type SessionSearchRequest struct {
	UserID     string `form:"user_id" binding:"omitempty,hexadecimal,len=24"`
	ActiveOnly *bool  `form:"active_only"` // true => no revocadas y no expiradas
}
