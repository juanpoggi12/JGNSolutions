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

type WorkoutEntryService struct {
	repo        *repositories.WorkoutEntryRepository
	sessionRepo *repositories.WorkoutSessionRepository
}

func NewWorkoutEntryService(repo *repositories.WorkoutEntryRepository, sessionRepo *repositories.WorkoutSessionRepository) *WorkoutEntryService {
	return &WorkoutEntryService{repo: repo, sessionRepo: sessionRepo}
}

// Crear una entrada de entrenamiento (verifica que la sesión pertenezca al actor)
func (s *WorkoutEntryService) Create(ctx context.Context, actor Actor, entry *models.WorkoutEntry) error {
	// Buscar la sesión a la que pertenece
	session, err := s.sessionRepo.FindByID(ctx, entry.WorkoutSessionID)
	if err != nil {
		return errors.New("la sesión de entrenamiento no existe")
	}

	// Validar permisos
	if actor.Role != "admin" && session.UserID != actor.UserID {
		return errors.New("no tienes permiso para agregar entradas a esta sesión")
	}

	// Asignar timestamps (opcional si tu modelo los tiene)
	_ = time.Now() // placeholder si luego agregás CreatedAt

	return s.repo.Create(ctx, entry)
}

// Obtener una entrada por ID (solo si pertenece a una sesión del actor o si es admin)
func (s *WorkoutEntryService) GetByID(ctx context.Context, actor Actor, id primitive.ObjectID) (*models.WorkoutEntry, error) {
	entry, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Traer la sesión a la que pertenece la entrada
	session, err := s.sessionRepo.FindByID(ctx, entry.WorkoutSessionID)
	if err != nil {
		return nil, errors.New("no se pudo validar la sesión asociada")
	}

	if actor.Role != "admin" && session.UserID != actor.UserID {
		return nil, errors.New("no tienes permiso para acceder a esta entrada")
	}

	return entry, nil
}

// Actualizar una entrada (solo si pertenece a una sesión del actor o si es admin)
func (s *WorkoutEntryService) Update(ctx context.Context, actor Actor, entry *models.WorkoutEntry) error {
	existing, err := s.repo.FindByID(ctx, entry.ID)
	if err != nil {
		return errors.New("entrada no encontrada")
	}

	session, err := s.sessionRepo.FindByID(ctx, existing.WorkoutSessionID)
	if err != nil {
		return errors.New("no se pudo validar la sesión asociada")
	}

	if actor.Role != "admin" && session.UserID != actor.UserID {
		return errors.New("no tienes permiso para modificar esta entrada")
	}

	return s.repo.Update(ctx, entry)
}

// Eliminar una entrada (solo si pertenece a una sesión del actor o si es admin)
func (s *WorkoutEntryService) Delete(ctx context.Context, actor Actor, id primitive.ObjectID) error {
	entry, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("entrada no encontrada")
	}

	session, err := s.sessionRepo.FindByID(ctx, entry.WorkoutSessionID)
	if err != nil {
		return errors.New("no se pudo validar la sesión asociada")
	}

	if actor.Role != "admin" && session.UserID != actor.UserID {
		return errors.New("no tienes permiso para eliminar esta entrada")
	}

	return s.repo.Delete(ctx, id)
}

// Buscar entradas (solo las que pertenecen a sesiones del actor, salvo que sea admin)
func (s *WorkoutEntryService) Search(ctx context.Context, actor Actor, filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutEntry, error) {
	if actor.Role != "admin" {
		// Filtramos por sesiones del usuario autenticado
		userSessions, err := s.sessionRepo.FindByUser(ctx, actor.UserID)
		if err != nil {
			return nil, errors.New("no se pudieron obtener las sesiones del usuario")
		}

		sessionIDs := make([]primitive.ObjectID, 0, len(userSessions))
		for _, s := range userSessions {
			sessionIDs = append(sessionIDs, s.ID)
		}

		filter["workoutSessionId"] = bson.M{"$in": sessionIDs}
	}

	return s.repo.Search(ctx, filter, opts...)
}
