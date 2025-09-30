package dto

type SessionResponse struct {
	ID         string  `json:"id"`
	UserID     string  `json:"user_id"`
	ExpiresAt  string  `json:"expires_at"`  // RFC3339 en el mapper
	CreatedAt  string  `json:"created_at"`  // RFC3339
	RevokedAt  *string `json:"revoked_at,omitempty"` // RFC3339
}

// Para revocar una sesiÃ³n especÃ­fica (por path: /sessions/:id/revoke)
type SessionIDUri struct {
	ID string `uri:"id" binding:"required,hexadecimal,len=24"`
}

// Para buscar sesiones de un usuario
type SessionSearchRequest struct {
	UserID      string `form:"user_id" binding:"omitempty,hexadecimal,len=24"`
	ActiveOnly  *bool  `form:"active_only"` // true => no revocadas y no expiradas
	PageQuery
}
