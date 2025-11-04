package services

import (
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogService struct {
	repository *repositories.LogRepository
}

func NewLogService(repo *repositories.LogRepository) *LogService {
	return &LogService{repository: repo}
}

// RecordAction → Registra una acción realizada por un usuario
func (s *LogService) RecordAction(userID primitive.ObjectID, action string) {
	log := models.Log{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now(),
	}

	_, _ = s.repository.InsertLog(log)
}
