package utils

import (
	"time"

    "github.com/juanpoggi12/JGNSolutions/backend/dto"
    "github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// CreateRequest → Model
func ConvertSessionCreateRequestToModel(req dto.SessionCreateRequest) (models.Session, error) {
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return models.Session{}, err
	}

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
		UserID:           userID,
		RefreshTokenHash: string(hashedToken),
		ExpiresAt:        expiresAt,
		CreatedAt:        time.Now(),
		RevokedAt:        nil,
	}, nil
}

// Model → Response
func ConvertSessionModelToResponse(s models.Session) dto.SessionResponse {
	var revokedAt *string
	if s.RevokedAt != nil {
		formatted := s.RevokedAt.Format(time.RFC3339)
		revokedAt = &formatted
	}

	return dto.SessionResponse{
		ID:        s.ID.Hex(),
		UserID:    s.UserID.Hex(),
		ExpiresAt: s.ExpiresAt.Format(time.RFC3339),
		CreatedAt: s.CreatedAt.Format(time.RFC3339),
		RevokedAt: revokedAt,
	}
}

// SearchRequest → filtro Mongo
func BuildSessionSearchFilter(search dto.SessionSearchRequest) bson.M {
	filter := bson.M{}
	if search.UserID != "" {
		filter["userId"] = ToObjectID(search.UserID)
	}
	if search.ActiveOnly != nil && *search.ActiveOnly {
		now := time.Now()
		filter["revokedAt"] = bson.M{"$exists": false}
		filter["expiresAt"] = bson.M{"$gt": now}
	}
	return filter
}
