package services

import (
	"errors"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"
	"github.com/juanpoggi12/JGNSolutions/backend/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExerciseServiceInterface interface {
	CreateExercise(actor Actor, req dto.ExerciseCreateRequest) (dto.ExerciseResponse, error)
	GetExerciseByID(actor Actor, id string) (dto.ExerciseResponse, error)
	UpdateExercise(actor Actor, id string, req dto.ExerciseUpdateRequest) (dto.ExerciseResponse, error)
	DeleteExercise(actor Actor, id string) error
	SearchExercises(actor Actor, search dto.ExerciseSearchRequest) ([]dto.ExerciseResponse, error)
}

type ExerciseService struct {
	repository repositories.ExerciseRepositoryInterface
	logService *LogService
}

func NewExerciseService(repository repositories.ExerciseRepositoryInterface, logService *LogService) *ExerciseService {
	return &ExerciseService{
		repository: repository,
		logService: logService,
	}
}

func (service *ExerciseService) CreateExercise(actor Actor, req dto.ExerciseCreateRequest) (dto.ExerciseResponse, error) {
	// Defensa adicional: solo admin puede crear ejercicios
	if actor.Role != "admin" {
		return dto.ExerciseResponse{}, errors.New("solo los administradores pueden crear ejercicios")
	}

	ejercicio, err := utils.ConvertExerciseCreateRequestToModel(req)
	if err != nil {
		return dto.ExerciseResponse{}, err
	}

	// Asignar auditoría
	ejercicio.CreatedBy = actor.UserID
	ejercicio.CreatedAt = time.Now()
	ejercicio.UpdatedAt = time.Now()

	resultado, err := service.repository.InsertarEjercicio(ejercicio)
	if err != nil {
		return dto.ExerciseResponse{}, err
	}

	if oid, ok := resultado.InsertedID.(primitive.ObjectID); ok {
		ejercicio.ID = oid

		if service.logService != nil {
			service.logService.RecordAction(actor.UserID, "ejercicio:creado "+ejercicio.Name)
		}

		return utils.ConvertExerciseModelToResponse(ejercicio), nil
	}

	return dto.ExerciseResponse{}, errors.New("error al obtener ID del ejercicio insertado")
}

func (service *ExerciseService) GetExerciseByID(actor Actor, id string) (dto.ExerciseResponse, error) {
	ejercicio, err := service.repository.ObtenerEjercicioPorID(id)
	if err != nil {
		return dto.ExerciseResponse{}, errors.New("ejercicio no encontrado")
	}
	if ejercicio.IsDeleted {
		return dto.ExerciseResponse{}, errors.New("el ejercicio fue eliminado")
	}

	return utils.ConvertExerciseModelToResponse(ejercicio), nil
}

func (service *ExerciseService) UpdateExercise(actor Actor, id string, req dto.ExerciseUpdateRequest) (dto.ExerciseResponse, error) {
	// Defensa adicional: solo admin puede modificar ejercicios
	if actor.Role != "admin" {
		return dto.ExerciseResponse{}, errors.New("solo los administradores pueden modificar ejercicios")
	}

	ejercicio, err := service.repository.ObtenerEjercicioPorID(id)
	if err != nil {
		return dto.ExerciseResponse{}, errors.New("ejercicio no encontrado")
	}
	if ejercicio.IsDeleted {
		return dto.ExerciseResponse{}, errors.New("el ejercicio fue eliminado")
	}

	utils.ApplyExerciseUpdateToModel(&ejercicio, req)
	ejercicio.UpdatedAt = time.Now()

	_, err = service.repository.ModificarEjercicio(ejercicio)
	if err != nil {
		return dto.ExerciseResponse{}, err
	}

	if service.logService != nil {
		service.logService.RecordAction(actor.UserID, "ejercicio:modificado "+ejercicio.Name)
	}

	return utils.ConvertExerciseModelToResponse(ejercicio), nil
}

func (service *ExerciseService) DeleteExercise(actor Actor, id string) error {
	// Defensa adicional: solo admin puede eliminar ejercicios
	if actor.Role != "admin" {
		return errors.New("solo los administradores pueden eliminar ejercicios")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID inválido")
	}

	resultado, err := service.repository.EliminarEjercicio(objectID)
	if err != nil {
		return err
	}
	if resultado.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	if service.logService != nil {
		service.logService.RecordAction(actor.UserID, "ejercicio:eliminado "+id)
	}

	return nil
}

func (service *ExerciseService) SearchExercises(actor Actor, search dto.ExerciseSearchRequest) ([]dto.ExerciseResponse, error) {
	// Los usuarios normales solo ven ejercicios no eliminados
	includeDeleted := false
	if actor.Role == "admin" {
		includeDeleted = search.IncludeDel
	}

	// Si es user, ignoramos CreatedBy del search (no debería poder filtrar por creador)
	createdBy := ""
	if actor.Role == "admin" && search.CreatedBy != "" {
		createdBy = search.CreatedBy
	}

	ejercicios, err := service.repository.BuscarEjercicios(
		search.Name,
		search.Category,
		search.MuscleGroup,
		search.Difficulty,
		createdBy,
		includeDeleted,
	)
	if err != nil {
		return nil, err
	}

	respuestas := make([]dto.ExerciseResponse, 0, len(ejercicios))
	for _, ejercicio := range ejercicios {
		respuestas = append(respuestas, utils.ConvertExerciseModelToResponse(ejercicio))
	}

	return respuestas, nil
}
