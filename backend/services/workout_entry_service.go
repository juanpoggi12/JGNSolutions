package services

import (
	"errors"
	"log"

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
type WorkoutEntryServiceInterface interface {
	Create(actor Actor, req dto.WorkoutEntryCreateRequest) (models.WorkoutEntry, error)
	Search(actor Actor, filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutEntry, error)
}

// Implementación
type WorkoutEntryService struct {
	repo        repositories.WorkoutEntryRepositoryInterface
	sessionRepo repositories.WorkoutSessionRepositoryInterface
	logService  *LogService
}

func NewWorkoutEntryService(repo repositories.WorkoutEntryRepositoryInterface, sessionRepo repositories.WorkoutSessionRepositoryInterface, logService *LogService) *WorkoutEntryService {
	return &WorkoutEntryService{repo: repo, sessionRepo: sessionRepo, logService: logService}
}

func (s *WorkoutEntryService) Create(actor Actor, req dto.WorkoutEntryCreateRequest) (models.WorkoutEntry, error) {
	sessionID, err := primitive.ObjectIDFromHex(req.WorkoutSessionID)
	if err != nil {
		return models.WorkoutEntry{}, errors.New("ID de sesión inválido")
	}

	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil {
		return models.WorkoutEntry{}, errors.New("la sesión de entrenamiento no existe")
	}

	if actor.Role != "admin" && session.UserID != actor.UserID {
		return models.WorkoutEntry{}, errors.New("no tienes permiso para agregar entradas a esta sesión")
	}

	entry, err := utils.ConvertWorkoutEntryCreateRequestToModel(req)
	if err != nil {
		return models.WorkoutEntry{}, err
	}

	if err := s.repo.Create(&entry); err != nil {
		return models.WorkoutEntry{}, err
	}

	if s.logService != nil {
		s.logService.RecordAction(actor.UserID, "entrada_entrenamiento:creada "+entry.ID.Hex())
	}

	return entry, nil
}

func (s *WorkoutEntryService) Search(actor Actor, filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutEntry, error) {
	log.Printf("[Service.Entry.Search] Actor: %s, Role: %s, Initial Filter: %v", actor.UserID.Hex(), actor.Role, filter)

	finalFilter := bson.M{}

	for k, v := range filter {
		finalFilter[k] = v
	}
	var specificSessionID primitive.ObjectID
	var isSpecificSessionSearch bool
	if sessionIDValue, ok := filter["workoutSessionId"]; ok {
		if oid, ok := sessionIDValue.(primitive.ObjectID); ok && !oid.IsZero() {
			specificSessionID = oid
			isSpecificSessionSearch = true
			log.Printf("[Service.Entry.Search] Detected specific session search for ID: %s", specificSessionID.Hex())
		} else {
			log.Printf("[Service.Entry.Search] Warning: workoutSessionId found in filter but is not a valid ObjectID: %v", sessionIDValue)

		}
	} else {
		log.Printf("[Service.Entry.Search] No specific workoutSessionId found in the initial filter.")
	}

	if actor.Role != "admin" {
		log.Printf("[Service.Entry.Search] Non-admin user detected.")

		if isSpecificSessionSearch {

			log.Printf("[Service.Entry.Search] Verifying ownership for session %s...", specificSessionID.Hex())
			requestedSession, err := s.sessionRepo.FindByID(specificSessionID)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					log.Printf("[Service.Entry.Search] Ownership check failed: Session %s not found.", specificSessionID.Hex())

					return []models.WorkoutEntry{}, nil
				}
				log.Printf("[Service.Entry.Search] Error checking session ownership for %s: %v", specificSessionID.Hex(), err)
				return nil, errors.New("error al verificar la sesión solicitada")
			}

			if requestedSession.UserID != actor.UserID {
				log.Printf("[Service.Entry.Search] Ownership check FAILED: Actor %s does not own session %s (Owner: %s).", actor.UserID.Hex(), specificSessionID.Hex(), requestedSession.UserID.Hex())

				return []models.WorkoutEntry{}, nil
			}

			log.Printf("[Service.Entry.Search] Ownership verified for session %s.", specificSessionID.Hex())

		} else {

			log.Printf("[Service.Entry.Search] Filtering entries by user's sessions (%s)...", actor.UserID.Hex())
			userSessions, err := s.sessionRepo.FindByUser(actor.UserID)
			if err != nil {
				log.Printf("[Service.Entry.Search] Error fetching user sessions: %v", err)
				return nil, errors.New("no se pudieron obtener las sesiones del usuario para filtrar entradas")
			}

			if len(userSessions) == 0 {
				log.Printf("[Service.Entry.Search] User has no sessions. Returning empty list.")
				return []models.WorkoutEntry{}, nil
			}

			sessionIDs := make([]primitive.ObjectID, 0, len(userSessions))
			for _, sess := range userSessions {
				sessionIDs = append(sessionIDs, sess.ID)
			}
			log.Printf("[Service.Entry.Search] Found %d sessions for user. Applying $in filter.", len(sessionIDs))

			finalFilter["workoutSessionId"] = bson.M{"$in": sessionIDs}
		}

	} else {

		log.Printf("[Service.Entry.Search] Admin user detected. No additional user filters applied.")

	}

	log.Printf("[Service.Entry.Search] Calling repository Search with FINAL filter: %v", finalFilter)
	results, err := s.repo.Search(finalFilter, opts...)
	if err != nil {
		log.Printf("[Service.Entry.Search] Error received from repository: %v", err)
		return nil, err
	}

	log.Printf("[Service.Entry.Search] Repository returned %d entries.", len(results))
	return results, nil
}
