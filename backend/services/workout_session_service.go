package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"
	"github.com/juanpoggi12/JGNSolutions/backend/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Interfaz del servicio
type WorkoutSessionServiceInterface interface {
	Create(actor Actor, req dto.WorkoutSessionCreateRequest) (models.WorkoutSession, error)
	GetByID(actor Actor, id primitive.ObjectID) (*models.WorkoutSession, error)
	Update(actor Actor, id primitive.ObjectID, req dto.WorkoutSessionUpdateRequest) (models.WorkoutSession, error)
	Delete(actor Actor, id primitive.ObjectID) error
	Search(actor Actor, filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutSession, error)
}

// Implementación
type WorkoutSessionService struct {
	repo       repositories.WorkoutSessionRepositoryInterface
	logService *LogService
}

func NewWorkoutSessionService(repo repositories.WorkoutSessionRepositoryInterface, logService *LogService) *WorkoutSessionService {
	return &WorkoutSessionService{repo: repo, logService: logService}
}

func (s *WorkoutSessionService) Create(actor Actor, req dto.WorkoutSessionCreateRequest) (models.WorkoutSession, error) {
	log.Printf("[Service.Create] Starting Create for actor %s", actor.UserID.Hex())
	if actor.Role != "user" && actor.Role != "admin" {
		log.Printf("[Service.Create] Authorization failed: Actor role '%s' not allowed.", actor.Role)
		return models.WorkoutSession{}, errors.New("rol no autorizado para crear sesiones")
	}
	session, err := utils.ConvertWorkoutSessionCreateRequestToModel(req)
	if err != nil {
		log.Printf("[Service.Create] Error converting DTO to model: %v", err)
		return models.WorkoutSession{}, fmt.Errorf("datos de sesión inválidos: %w", err)
	}

	session.UserID = actor.UserID
	now := time.Now()
	if session.CreatedAt.IsZero() {
		session.CreatedAt = now
	}
	session.UpdatedAt = now

	log.Printf("[Service.Create] Model prepared for repo. UserID: %s, RoutineID: %v", session.UserID.Hex(), session.RoutineID)

	if err := s.repo.Create(&session); err != nil {
		log.Printf("[Service.Create] Error received from repo.Create: %v", err)

		return models.WorkoutSession{}, err
	}

	if session.ID.IsZero() {

		log.Printf("[Service.Create] CRITICAL ERROR: repo.Create returned nil error, but session.ID is still Zero!")
		return models.WorkoutSession{}, errors.New("error crítico: no se pudo asignar ID a la sesión")
	}
	log.Printf("[Service.Create] repo.Create succeeded. Session object ID in service is now: %s", session.ID.Hex())

	if s.logService != nil {
		s.logService.RecordAction(actor.UserID, "sesion_entrenamiento:creada "+session.ID.Hex())
	}

	log.Printf("[Service.Create] Returning successful session model with ID: %s", session.ID.Hex())
	return session, nil
}

// Obtener una sesión por ID
func (s *WorkoutSessionService) GetByID(actor Actor, id primitive.ObjectID) (*models.WorkoutSession, error) {
	log.Printf("[Service.GetByID] Getting session %s for actor %s", id.Hex(), actor.UserID.Hex())
	session, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("[Service.GetByID] Session %s not found.", id.Hex())
			return nil, errors.New("sesión no encontrada")
		}
		log.Printf("[Service.GetByID] Error finding session %s: %v", id.Hex(), err)
		return nil, err
	}

	if actor.Role != "admin" && session.UserID != actor.UserID {
		log.Printf("[Service.GetByID] Authorization failed: Actor %s cannot access session %s owned by %s.", actor.UserID.Hex(), id.Hex(), session.UserID.Hex())
		return nil, errors.New("no tienes permiso para acceder a esta sesión")
	}

	log.Printf("[Service.GetByID] Access granted for session %s.", id.Hex())
	return session, nil
}

// Actualizar una sesión
func (s *WorkoutSessionService) Update(actor Actor, id primitive.ObjectID, req dto.WorkoutSessionUpdateRequest) (models.WorkoutSession, error) {
	log.Printf("[Service.Update] Updating session %s for actor %s", id.Hex(), actor.UserID.Hex())
	existing, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("[Service.Update] Session %s not found for update.", id.Hex())
			return models.WorkoutSession{}, errors.New("sesión no encontrada para actualizar")
		}
		log.Printf("[Service.Update] Error finding session %s for update: %v", id.Hex(), err)
		return models.WorkoutSession{}, err
	}

	if actor.Role != "admin" && existing.UserID != actor.UserID {
		log.Printf("[Service.Update] Authorization failed: Actor %s cannot update session %s owned by %s.", actor.UserID.Hex(), id.Hex(), existing.UserID.Hex())
		return models.WorkoutSession{}, errors.New("no tienes permiso para modificar esta sesión")
	}

	if err := utils.ApplyWorkoutSessionUpdateToModel(existing, req); err != nil {
		log.Printf("[Service.Update] Error applying updates to model for session %s: %v", id.Hex(), err)
		return models.WorkoutSession{}, fmt.Errorf("error en los datos de actualización: %w", err)
	}
	existing.UpdatedAt = time.Now()

	log.Printf("[Service.Update] Model updated, calling repo.Update for session %s.", id.Hex())
	if err := s.repo.Update(existing); err != nil {
		log.Printf("[Service.Update] Error received from repo.Update for session %s: %v", id.Hex(), err)
		return models.WorkoutSession{}, err
	}

	if s.logService != nil {
		s.logService.RecordAction(actor.UserID, "sesion_entrenamiento:modificada "+existing.ID.Hex())
	}

	log.Printf("[Service.Update] Successfully updated session %s.", id.Hex())
	return *existing, nil
}

func (s *WorkoutSessionService) Delete(actor Actor, id primitive.ObjectID) error {
	log.Printf("[Service.Delete] Deleting session %s for actor %s", id.Hex(), actor.UserID.Hex())

	session, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("[Service.Delete] Session %s not found for deletion.", id.Hex())
			return errors.New("sesión no encontrada para eliminar")
		}
		log.Printf("[Service.Delete] Error finding session %s for deletion: %v", id.Hex(), err)
		return err
	}

	if actor.Role != "admin" && session.UserID != actor.UserID {
		log.Printf("[Service.Delete] Authorization failed: Actor %s cannot delete session %s owned by %s.", actor.UserID.Hex(), id.Hex(), session.UserID.Hex())
		return errors.New("no tienes permiso para eliminar esta sesión")
	}

	log.Printf("[Service.Delete] Permissions ok, calling repo.Delete for session %s.", id.Hex())
	if err := s.repo.Delete(id); err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {

			log.Printf("[Service.Delete] Session %s disappeared before deletion?", id.Hex())
			return errors.New("la sesión no se encontró al intentar eliminar")
		}
		return err
	}

	if s.logService != nil {
		s.logService.RecordAction(actor.UserID, "sesion_entrenamiento:eliminada "+id.Hex())
	}

	log.Printf("[Service.Delete] Successfully deleted session %s.", id.Hex())
	return nil
}

func (s *WorkoutSessionService) Search(actor Actor, filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutSession, error) {
	log.Printf("[Service.Search] Searching sessions for actor %s with filter: %v", actor.UserID.Hex(), filter)

	if actor.Role != "admin" {
		log.Printf("[Service.Search] Applying user filter for non-admin user %s", actor.UserID.Hex())
		filter["userId"] = actor.UserID
	} else {
		log.Printf("[Service.Search] Admin user %s requested search, no user filter applied.", actor.UserID.Hex())
	}

	return s.repo.Search(filter, opts...)
}
