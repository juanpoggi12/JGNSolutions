// juanpoggi12/jgnsolutions/JGNSolutions-7c48b53190321ccfabc6877d44ae535f756457c5/backend/services/workout_session_service.go
package services

import (
	"errors"
	"fmt" // Import fmt package
	"log" // Import log package
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"
	"github.com/juanpoggi12/JGNSolutions/backend/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo" // Import mongo package for ErrNoDocuments
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

// Crear una nueva sesión
func (s *WorkoutSessionService) Create(actor Actor, req dto.WorkoutSessionCreateRequest) (models.WorkoutSession, error) {
	log.Printf("[Service.Create] Starting Create for actor %s", actor.UserID.Hex())
	if actor.Role != "user" && actor.Role != "admin" {
		log.Printf("[Service.Create] Authorization failed: Actor role '%s' not allowed.", actor.Role)
		return models.WorkoutSession{}, errors.New("rol no autorizado para crear sesiones")
	}

	// Convert DTO, sets times internally
	session, err := utils.ConvertWorkoutSessionCreateRequestToModel(req)
	if err != nil {
		log.Printf("[Service.Create] Error converting DTO to model: %v", err)
		return models.WorkoutSession{}, fmt.Errorf("datos de sesión inválidos: %w", err)
	}

	// Assign owner and ensure timestamps are set
	session.UserID = actor.UserID
	now := time.Now()
	if session.CreatedAt.IsZero() { // Should be set by util, but double check
		session.CreatedAt = now
	}
	session.UpdatedAt = now // Always set UpdatedAt on creation

	log.Printf("[Service.Create] Model prepared for repo. UserID: %s, RoutineID: %v", session.UserID.Hex(), session.RoutineID)

	// Call repository Create, passing the pointer
	// The session object *should* be updated with the ID by the repo function
	if err := s.repo.Create(&session); err != nil {
		log.Printf("[Service.Create] Error received from repo.Create: %v", err)
		// Propagate the error directly
		return models.WorkoutSession{}, err
	}

	// --- Log ID Verification ---
	// Check the ID *after* the repo call returns
	if session.ID.IsZero() {
		// This indicates the repo failed to assign a valid ID despite returning nil error
		log.Printf("[Service.Create] CRITICAL ERROR: repo.Create returned nil error, but session.ID is still Zero!")
		return models.WorkoutSession{}, errors.New("error crítico: no se pudo asignar ID a la sesión")
	}
	log.Printf("[Service.Create] repo.Create succeeded. Session object ID in service is now: %s", session.ID.Hex())
	// --- End Log ID Verification ---

	// Log the action using the (now confirmed) session ID
	if s.logService != nil {
		s.logService.RecordAction(actor.UserID, "sesion_entrenamiento:creada "+session.ID.Hex())
	}

	// Return the session model (which now includes the ID)
	log.Printf("[Service.Create] Returning successful session model with ID: %s", session.ID.Hex())
	return session, nil
}

// Obtener una sesión por ID
func (s *WorkoutSessionService) GetByID(actor Actor, id primitive.ObjectID) (*models.WorkoutSession, error) {
	log.Printf("[Service.GetByID] Getting session %s for actor %s", id.Hex(), actor.UserID.Hex())
	session, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) { // Usa errors.Is para comparar errores estándar
			log.Printf("[Service.GetByID] Session %s not found.", id.Hex())
			return nil, errors.New("sesión no encontrada")
		}
		log.Printf("[Service.GetByID] Error finding session %s: %v", id.Hex(), err)
		return nil, err // Otro error
	}

	// Check permissions
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

	// Check permissions
	if actor.Role != "admin" && existing.UserID != actor.UserID {
		log.Printf("[Service.Update] Authorization failed: Actor %s cannot update session %s owned by %s.", actor.UserID.Hex(), id.Hex(), existing.UserID.Hex())
		return models.WorkoutSession{}, errors.New("no tienes permiso para modificar esta sesión")
	}

	// Apply updates using the utility function (or manually as before)
	if err := utils.ApplyWorkoutSessionUpdateToModel(existing, req); err != nil {
		log.Printf("[Service.Update] Error applying updates to model for session %s: %v", id.Hex(), err)
		return models.WorkoutSession{}, fmt.Errorf("error en los datos de actualización: %w", err)
	}
	existing.UpdatedAt = time.Now() // Ensure UpdatedAt is set

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

// Eliminar una sesión
func (s *WorkoutSessionService) Delete(actor Actor, id primitive.ObjectID) error {
	log.Printf("[Service.Delete] Deleting session %s for actor %s", id.Hex(), actor.UserID.Hex())
	// First, check if session exists and if user has permission
	session, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("[Service.Delete] Session %s not found for deletion.", id.Hex())
			return errors.New("sesión no encontrada para eliminar")
		}
		log.Printf("[Service.Delete] Error finding session %s for deletion: %v", id.Hex(), err)
		return err
	}

	// Check permissions
	if actor.Role != "admin" && session.UserID != actor.UserID {
		log.Printf("[Service.Delete] Authorization failed: Actor %s cannot delete session %s owned by %s.", actor.UserID.Hex(), id.Hex(), session.UserID.Hex())
		return errors.New("no tienes permiso para eliminar esta sesión")
	}

	// Proceed with deletion
	log.Printf("[Service.Delete] Permissions ok, calling repo.Delete for session %s.", id.Hex())
	if err := s.repo.Delete(id); err != nil {
		// repo.Delete already logs the specific error
		if errors.Is(err, mongo.ErrNoDocuments) {
			// Should not happen if FindByID succeeded, but handle defensively
			log.Printf("[Service.Delete] Session %s disappeared before deletion?", id.Hex())
			return errors.New("la sesión no se encontró al intentar eliminar")
		}
		return err // Other repo error
	}

	// TODO: Consider deleting associated WorkoutEntry documents here or using transactions/DB features

	if s.logService != nil {
		s.logService.RecordAction(actor.UserID, "sesion_entrenamiento:eliminada "+id.Hex())
	}

	log.Printf("[Service.Delete] Successfully deleted session %s.", id.Hex())
	return nil
}

// Buscar sesiones (no se registran logs)
func (s *WorkoutSessionService) Search(actor Actor, filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutSession, error) {
	log.Printf("[Service.Search] Searching sessions for actor %s with filter: %v", actor.UserID.Hex(), filter)
	// Apply user filter if not admin
	if actor.Role != "admin" {
		log.Printf("[Service.Search] Applying user filter for non-admin user %s", actor.UserID.Hex())
		filter["userId"] = actor.UserID
	} else {
		log.Printf("[Service.Search] Admin user %s requested search, no user filter applied.", actor.UserID.Hex())
	}
	// The repository handles default sorting if not provided in opts
	return s.repo.Search(filter, opts...)
}
