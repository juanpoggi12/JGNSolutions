package utils

import (
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// CreateRequest → Model
// Ya no toma el UserID del DTO, se asigna externamente en el service usando actor.UserID
func ConvertSessionCreateRequestToModel(req dto.SessionCreateRequest) (models.Session, error) {
	// Parseo de fecha (RFC3339)
	expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt)
	if err != nil {
		return models.Session{}, err
	}

	// Hashear el refresh token antes de guardarlo
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(req.Token), bcrypt.DefaultCost)
	if err != nil {
		return models.Session{}, err
	}

	return models.Session{
		// UserID se asignará en el service
		RefreshTokenHash: string(hashedToken),
		ExpiresAt:        expiresAt,
		CreatedAt:        time.Now(),
		RevokedAt:        nil,
	}, nil
}

// Model → Response
// Quita el campo UserID porque ya no está en el DTO
func ConvertSessionModelToResponse(s models.Session) dto.SessionResponse {
	var revokedAt *string
	if s.RevokedAt != nil {
		formatted := s.RevokedAt.Format(time.RFC3339)
		revokedAt = &formatted
	}

	return dto.SessionResponse{
		ID:        s.ID.Hex(),
		ExpiresAt: s.ExpiresAt.Format(time.RFC3339),
		CreatedAt: s.CreatedAt.Format(time.RFC3339),
		RevokedAt: revokedAt,
	}
}

// SearchRequest → filtro Mongo
// Ya no usa UserID del body; el filtro por usuario lo manejará el service usando actor.UserID
func BuildSessionSearchFilter(search dto.SessionSearchRequest, userID string) bson.M {
	filter := bson.M{}

	if userID != "" {
		filter["userId"] = ToObjectID(userID)
	}
	if search.ActiveOnly != nil && *search.ActiveOnly {
		now := time.Now()
		filter["revokedAt"] = bson.M{"$exists": false}
		filter["expiresAt"] = bson.M{"$gt": now}
	}

	return filter
}
