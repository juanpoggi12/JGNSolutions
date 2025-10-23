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

// 🧩 Interfaz del servicio
type WorkoutEntryServiceInterface interface {
	Create(actor Actor, req dto.WorkoutEntryCreateRequest) (models.WorkoutEntry, error)
	GetByID(actor Actor, id primitive.ObjectID) (*models.WorkoutEntry, error)
	Update(actor Actor, id primitive.ObjectID, req dto.WorkoutEntryUpdateRequest) (models.WorkoutEntry, error)
	Delete(actor Actor, id primitive.ObjectID) error
	Search(actor Actor, filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutEntry, error)
}

// 🧩 Implementación
type WorkoutEntryService struct {
	repo        repositories.WorkoutEntryRepositoryInterface
	sessionRepo repositories.WorkoutSessionRepositoryInterface
	logService  *LogService
}

func NewWorkoutEntryService(repo repositories.WorkoutEntryRepositoryInterface, sessionRepo repositories.WorkoutSessionRepositoryInterface, logService *LogService) *WorkoutEntryService {
	return &WorkoutEntryService{repo: repo, sessionRepo: sessionRepo, logService: logService}
}

// Crear una nueva entrada
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

// Obtener una entrada por ID
func (s *WorkoutEntryService) GetByID(actor Actor, id primitive.ObjectID) (*models.WorkoutEntry, error) {
	entry, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("entrada no encontrada")
	}

	session, err := s.sessionRepo.FindByID(entry.WorkoutSessionID)
	if err != nil {
		return nil, errors.New("no se pudo validar la sesión asociada")
	}

	if actor.Role != "admin" && session.UserID != actor.UserID {
		return nil, errors.New("no tienes permiso para acceder a esta entrada")
	}

	return entry, nil
}

// Actualizar una entrada
func (s *WorkoutEntryService) Update(actor Actor, id primitive.ObjectID, req dto.WorkoutEntryUpdateRequest) (models.WorkoutEntry, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return models.WorkoutEntry{}, errors.New("entrada no encontrada")
	}

	session, err := s.sessionRepo.FindByID(existing.WorkoutSessionID)
	if err != nil {
		return models.WorkoutEntry{}, errors.New("no se pudo validar la sesión asociada")
	}

	if actor.Role != "admin" && session.UserID != actor.UserID {
		return models.WorkoutEntry{}, errors.New("no tienes permiso para modificar esta entrada")
	}

	// Aplicar campos enviados
	if req.SetNumber != nil {
		existing.SetNumber = *req.SetNumber
	}
	if req.RepsDone != nil {
		existing.RepsDone = req.RepsDone
	}
	if req.WeightUsed != nil {
		existing.WeightUsed = req.WeightUsed
	}
	if req.TimeSec != nil {
		existing.TimeSec = req.TimeSec
	}
	if req.PerceivedEffort != nil {
		existing.PerceivedEffort = req.PerceivedEffort
	}

	if err := s.repo.Update(existing); err != nil {
		return models.WorkoutEntry{}, err
	}

	if s.logService != nil {
		s.logService.RecordAction(actor.UserID, "entrada_entrenamiento:modificada "+existing.ID.Hex())
	}

	return *existing, nil
}

// Eliminar una entrada
func (s *WorkoutEntryService) Delete(actor Actor, id primitive.ObjectID) error {
	entry, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("entrada no encontrada")
	}

	session, err := s.sessionRepo.FindByID(entry.WorkoutSessionID)
	if err != nil {
		return errors.New("no se pudo validar la sesión asociada")
	}

	if actor.Role != "admin" && session.UserID != actor.UserID {
		return errors.New("no tienes permiso para eliminar esta entrada")
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	if s.logService != nil {
		s.logService.RecordAction(actor.UserID, "entrada_entrenamiento:eliminada "+id.Hex())
	}

	return nil
}

// juanpoggi12/jgnsolutions/JGNSolutions-7c48b53190321ccfabc6877d44ae535f756457c5/backend/services/workout_entry_service.go

// Buscar entradas (no se registran logs)
func (s *WorkoutEntryService) Search(actor Actor, filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutEntry, error) {
	log.Printf("[Service.Entry.Search] Actor: %s, Initial Filter: %v", actor.UserID.Hex(), filter) // Log inicial

	// Verifica si ya se está filtrando por una sesión específica
	_, specificSessionSearch := filter["workoutSessionId"]

	// Si NO es admin Y NO se está buscando ya una sesión específica,
	// entonces restringe la búsqueda a las sesiones del usuario.
	if actor.Role != "admin" && !specificSessionSearch {
		log.Printf("[Service.Entry.Search] Non-admin user, no specific session requested. Filtering by user's sessions.")
		userSessions, err := s.sessionRepo.FindByUser(actor.UserID)
		if err != nil {
			log.Printf("[Service.Entry.Search] Error fetching user sessions: %v", err)
			return nil, errors.New("no se pudieron obtener las sesiones del usuario para filtrar entradas")
		}

		if len(userSessions) == 0 {
			log.Printf("[Service.Entry.Search] User has no sessions. Returning empty list.")
			return []models.WorkoutEntry{}, nil // Si el usuario no tiene sesiones, no puede tener entradas
		}

		sessionIDs := make([]primitive.ObjectID, 0, len(userSessions))
		for _, sess := range userSessions {
			sessionIDs = append(sessionIDs, sess.ID)
		}
		log.Printf("[Service.Entry.Search] Applying filter for user's sessions: %v", sessionIDs)
		filter["workoutSessionId"] = bson.M{"$in": sessionIDs} // Aplica filtro por TODAS las sesiones del usuario

	} else if actor.Role != "admin" && specificSessionSearch {
		// Si NO es admin PERO SÍ busca una sesión específica, debemos verificar que esa sesión le pertenezca.
		log.Printf("[Service.Entry.Search] Non-admin user requested specific session. Verifying ownership.")
		sessionID, ok := filter["workoutSessionId"].(primitive.ObjectID)
		if !ok {
			// Si el filtro no es un ObjectID (raro, pero posible si viene mal del handler)
			log.Printf("[Service.Entry.Search] Error: workoutSessionId filter is not a primitive.ObjectID. Filter value: %v", filter["workoutSessionId"])
			return nil, errors.New("filtro de sesión inválido")
		}

		// Busca la sesión específica para verificar el propietario
		requestedSession, err := s.sessionRepo.FindByID(sessionID)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				log.Printf("[Service.Entry.Search] Ownership check failed: Session %s not found.", sessionID.Hex())
				return nil, errors.New("sesión solicitada no encontrada")
			}
			log.Printf("[Service.Entry.Search] Error checking session ownership for %s: %v", sessionID.Hex(), err)
			return nil, errors.New("error al verificar la sesión solicitada")
		}

		// Compara el UserID de la sesión con el del actor
		if requestedSession.UserID != actor.UserID {
			log.Printf("[Service.Entry.Search] Ownership check failed: Actor %s does not own session %s (Owner: %s).", actor.UserID.Hex(), sessionID.Hex(), requestedSession.UserID.Hex())
			return nil, errors.New("no tienes permiso para acceder a las entradas de esta sesión")
		}
		log.Printf("[Service.Entry.Search] Ownership verified for session %s.", sessionID.Hex())
		// Si la verificación pasa, el filtro original (filter["workoutSessionId"] = sessionID) se mantiene y se usa.

	} else {
		// Si es admin, puede buscar cualquier sesión o todas.
		log.Printf("[Service.Entry.Search] Admin user. Proceeding with filter: %v", filter)
	}

	// Llama al repositorio con el filtro final (modificado o no)
	log.Printf("[Service.Entry.Search] Calling repository Search with final filter: %v", filter)
	return s.repo.Search(filter, opts...)
}
