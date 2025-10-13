package services

import (
	"errors"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"
	"github.com/juanpoggi12/JGNSolutions/backend/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 🧩 Interfaz del servicio
type WorkoutSessionServiceInterface interface {
	Create(actor Actor, req dto.WorkoutSessionCreateRequest) (models.WorkoutSession, error)
	GetByID(actor Actor, id primitive.ObjectID) (*models.WorkoutSession, error)
	Update(actor Actor, id primitive.ObjectID, req dto.WorkoutSessionUpdateRequest) (models.WorkoutSession, error)
	Delete(actor Actor, id primitive.ObjectID) error
	Search(actor Actor, filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutSession, error)
}

// 🧩 Implementación
type WorkoutSessionService struct {
	repo repositories.WorkoutSessionRepositoryInterface
}

func NewWorkoutSessionService(repo repositories.WorkoutSessionRepositoryInterface) *WorkoutSessionService {
	return &WorkoutSessionService{repo: repo}
}

// Crear una nueva sesión
func (s *WorkoutSessionService) Create(actor Actor, req dto.WorkoutSessionCreateRequest) (models.WorkoutSession, error) {
	if actor.Role != "user" && actor.Role != "admin" {
		return models.WorkoutSession{}, errors.New("rol no autorizado para crear sesiones")
	}

	session, err := utils.ConvertWorkoutSessionCreateRequestToModel(req)
	if err != nil {
		return models.WorkoutSession{}, err
	}

	session.UserID = actor.UserID
	session.CreatedAt = time.Now()

	if err := s.repo.Create(&session); err != nil {
		return models.WorkoutSession{}, err
	}

	return session, nil
}

// Obtener una sesión por ID
func (s *WorkoutSessionService) GetByID(actor Actor, id primitive.ObjectID) (*models.WorkoutSession, error) {
	session, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("sesión no encontrada")
	}

	if actor.Role != "admin" && session.UserID != actor.UserID {
		return nil, errors.New("no tienes permiso para acceder a esta sesión")
	}

	return session, nil
}

// Actualizar una sesión
func (s *WorkoutSessionService) Update(actor Actor, id primitive.ObjectID, req dto.WorkoutSessionUpdateRequest) (models.WorkoutSession, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return models.WorkoutSession{}, errors.New("sesión no encontrada")
	}

	if actor.Role != "admin" && existing.UserID != actor.UserID {
		return models.WorkoutSession{}, errors.New("no tienes permiso para modificar esta sesión")
	}

	// Aplicar solo los campos enviados
	if req.RoutineID != nil {
		if oid, err := primitive.ObjectIDFromHex(*req.RoutineID); err == nil {
			existing.RoutineID = &oid
		}
	}
	if req.StartTime != nil {
		start, _ := time.Parse(time.RFC3339, *req.StartTime)
		existing.StartTime = start
	}
	if req.EndTime != nil {
		end, _ := time.Parse(time.RFC3339, *req.EndTime)
		existing.EndTime = end
	}
	if req.Notes != nil {
		existing.Notes = *req.Notes
	}
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(existing); err != nil {
		return models.WorkoutSession{}, err
	}

	return *existing, nil
}

// Eliminar una sesión
func (s *WorkoutSessionService) Delete(actor Actor, id primitive.ObjectID) error {
	session, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("sesión no encontrada")
	}

	if actor.Role != "admin" && session.UserID != actor.UserID {
		return errors.New("no tienes permiso para eliminar esta sesión")
	}

	return s.repo.Delete(id)
}

// Buscar sesiones (solo las del actor, salvo admin)
func (s *WorkoutSessionService) Search(actor Actor, filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutSession, error) {
	if actor.Role != "admin" {
		filter["userId"] = actor.UserID
	}
	return s.repo.Search(filter, opts...)
}
