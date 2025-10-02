package utils

import (
	"JGNSolutions/dto"
	"JGNSolutions/models"
)

func ConvertSessionRequestToModel(req dto.SessionRequest) models.Session {
	return models.Session{
		UserID:    ToObjectID(req.UserID),
		Token:     req.Token,
		ExpiresAt: req.ExpiresAt,
	}
}

func ConvertSessionModelToResponse(s models.Session) dto.SessionResponse {
	return dto.SessionResponse{
		ID:        s.ID.Hex(),
		UserID:    s.UserID.Hex(),
		Token:     s.Token,
		ExpiresAt: s.ExpiresAt,
	}
}
