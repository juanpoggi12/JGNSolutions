package services

import (
	"context"
	"errors"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WorkoutSessionService struct {
	repo *repositories.WorkoutSessionRepository
}

func NewWorkoutSessionService(repo *repositories.WorkoutSessionRepository) *WorkoutSessionService {
	return &WorkoutSessionService{repo: repo}
}

// Crear una nueva sesión de entrenamiento para el usuario autenticado
func (s *WorkoutSessionService) Create(ctx context.Context, actor Actor, session *models.WorkoutSession) error {
	if actor.Role != "user" && actor.Role != "admin" {
		return errors.New("rol no autorizado para crear sesiones")
	}

	session.UserID = actor.UserID
	session.CreatedAt = time.Now()

	return s.repo.Create(ctx, session)
}

// Obtener una sesión por ID (solo si pertenece al actor o si es admin)
func (s *WorkoutSessionService) GetByID(ctx context.Context, actor Actor, id primitive.ObjectID) (*models.WorkoutSession, error) {
	session, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Los usuarios solo pueden ver sus propias sesiones
	if actor.Role != "admin" && session.UserID != actor.UserID {
		return nil, errors.New("no tienes permiso para acceder a esta sesión")
	}

	return session, nil
}

// Actualizar una sesión (solo si pertenece al actor o si es admin)
func (s *WorkoutSessionService) Update(ctx context.Context, actor Actor, session *models.WorkoutSession) error {
	existing, err := s.repo.FindByID(ctx, session.ID)
	if err != nil {
		return errors.New("sesión no encontrada")
	}

	if actor.Role != "admin" && existing.UserID != actor.UserID {
		return errors.New("no tienes permiso para modificar esta sesión")
	}

	session.UserID = existing.UserID // evitar cambiar el propietario
	return s.repo.Update(ctx, session)
}

// Eliminar una sesión (solo si pertenece al actor o si es admin)
func (s *WorkoutSessionService) Delete(ctx context.Context, actor Actor, id primitive.ObjectID) error {
	session, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("sesión no encontrada")
	}

	if actor.Role != "admin" && session.UserID != actor.UserID {
		return errors.New("no tienes permiso para eliminar esta sesión")
	}

	return s.repo.Delete(ctx, id)
}

// Buscar sesiones (solo las del actor, salvo que sea admin)
func (s *WorkoutSessionService) Search(ctx context.Context, actor Actor, filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutSession, error) {
	if actor.Role != "admin" {
		filter["userId"] = actor.UserID
	}
	return s.repo.Search(ctx, filter, opts...)
}
