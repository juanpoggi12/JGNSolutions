package services

import (
	"context"
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

// Registrar una acción del usuario
func (s *LogService) RecordAction(ctx context.Context, userID primitive.ObjectID, action string) {
	log := models.Log{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now(),
	}
	_, _ = s.repository.InsertarLog(ctx, &log) // Ignora errores, no interrumpe flujo
}

// Listar todos los logs
func (s *LogService) ListarLogs(ctx context.Context) ([]models.Log, error) {
	return s.repository.ListarLogs(ctx)
}
