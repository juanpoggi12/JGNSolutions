package services

import (
	"context"
	"errors"
	"strings"
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
	ListExercises(actor Actor, query dto.ExerciseCatalogQuery) (dto.ExerciseCatalogResponse, error)
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

	exercise, err := utils.ConvertExerciseCreateRequestToModel(req)
	if err != nil {
		return dto.ExerciseResponse{}, err
	}

	// Asignar auditoría
	exercise.CreatedBy = actor.UserID
	exercise.CreatedAt = time.Now()
	exercise.UpdatedAt = time.Now()

	result, err := service.repository.InsertExercise(exercise)
	if err != nil {
		return dto.ExerciseResponse{}, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		exercise.ID = oid

		if service.logService != nil {
			service.logService.RecordAction(actor.UserID, "exercise:created "+exercise.Name)
		}

		return utils.ConvertExerciseModelToResponse(exercise), nil
	}

	return dto.ExerciseResponse{}, errors.New("error al obtener ID del ejercicio insertado")
}

func (service *ExerciseService) GetExerciseByID(actor Actor, id string) (dto.ExerciseResponse, error) {
	ejercicio, err := service.repository.GetExerciseByID(id)
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

	ejercicio, err := service.repository.GetExerciseByID(id)
	if err != nil {
		return dto.ExerciseResponse{}, errors.New("ejercicio no encontrado")
	}
	if ejercicio.IsDeleted {
		return dto.ExerciseResponse{}, errors.New("el ejercicio fue eliminado")
	}

	utils.ApplyExerciseUpdateToModel(&ejercicio, req)
	ejercicio.UpdatedAt = time.Now()

	_, err = service.repository.UpdateExercise(ejercicio)
	if err != nil {
		return dto.ExerciseResponse{}, err
	}

	if service.logService != nil {
		service.logService.RecordAction(actor.UserID, "exercise:updated "+ejercicio.Name)
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

	resultado, err := service.repository.DeleteExercise(objectID)
	if err != nil {
		return err
	}
	if resultado.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	if service.logService != nil {
		service.logService.RecordAction(actor.UserID, "exercise:deleted "+id)
	}

	return nil
}

func (service *ExerciseService) ListExercises(actor Actor, query dto.ExerciseCatalogQuery) (dto.ExerciseCatalogResponse, error) {
	page := query.Page
	if page <= 0 {
		page = 1
	}

	limit := query.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	params := repositories.ExerciseCatalogParams{
		Query:          strings.TrimSpace(query.Q),
		Category:       strings.TrimSpace(query.Category),
		MuscleGroup:    strings.TrimSpace(query.MuscleGroup),
		Difficulty:     strings.TrimSpace(query.Difficulty),
		IncludeDeleted: actor.Role == "admin",
		Page:           page,
		Limit:          limit,
	}

	ctx := context.Background()
	ejercicios, total, err := service.repository.ListExercisesCatalog(ctx, params)
	if err != nil {
		return dto.ExerciseCatalogResponse{}, err
	}

	items := make([]dto.ExerciseResponse, 0, len(ejercicios))
	for _, ejercicio := range ejercicios {
		items = append(items, utils.ConvertExerciseModelToResponse(ejercicio))
	}

	return dto.ExerciseCatalogResponse{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}
