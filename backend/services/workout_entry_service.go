package services

import (
	"errors"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"
	"github.com/juanpoggi12/JGNSolutions/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Buscar entradas (no se registran logs)
func (s *WorkoutEntryService) Search(actor Actor, filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutEntry, error) {
	if actor.Role != "admin" {
		userSessions, err := s.sessionRepo.FindByUser(actor.UserID)
		if err != nil {
			return nil, errors.New("no se pudieron obtener las sesiones del usuario")
		}

		sessionIDs := make([]primitive.ObjectID, 0, len(userSessions))
		for _, s := range userSessions {
			sessionIDs = append(sessionIDs, s.ID)
		}

		filter["workoutSessionId"] = bson.M{"$in": sessionIDs}
	}

	return s.repo.Search(filter, opts...)
}
