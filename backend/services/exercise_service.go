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
	CreateExercise(req dto.ExerciseCreateRequest) (dto.ExerciseResponse, error)
	GetExerciseByID(id string) (dto.ExerciseResponse, error)
	UpdateExercise(id string, req dto.ExerciseUpdateRequest) (dto.ExerciseResponse, error)
	DeleteExercise(id string) error
	SearchExercises(search dto.ExerciseSearchRequest) ([]dto.ExerciseResponse, error)
}

type ExerciseService struct {
	repository repositories.ExerciseRepositoryInterface
}

func NewExerciseService(repository repositories.ExerciseRepositoryInterface) *ExerciseService {
	return &ExerciseService{repository: repository}
}

func (service *ExerciseService) CreateExercise(req dto.ExerciseCreateRequest) (dto.ExerciseResponse, error) {
	ejercicio, err := utils.ConvertExerciseCreateRequestToModel(req)
	if err != nil {
		return dto.ExerciseResponse{}, err
	}

	resultado, err := service.repository.InsertarEjercicio(ejercicio)
	if err != nil {
		return dto.ExerciseResponse{}, err
	}

	if oid, ok := resultado.InsertedID.(primitive.ObjectID); ok {
		ejercicio.ID = oid
		return utils.ConvertExerciseModelToResponse(ejercicio), nil
	}

	return dto.ExerciseResponse{}, errors.New("error al obtener ID del ejercicio insertado")
}

func (service *ExerciseService) GetExerciseByID(id string) (dto.ExerciseResponse, error) {
	ejercicio, err := service.repository.ObtenerEjercicioPorID(id)
	if err != nil {
		return dto.ExerciseResponse{}, errors.New("ejercicio no encontrado")
	}
	if ejercicio.IsDeleted {
		return dto.ExerciseResponse{}, errors.New("el ejercicio fue eliminado")
	}

	return utils.ConvertExerciseModelToResponse(ejercicio), nil
}

func (service *ExerciseService) UpdateExercise(id string, req dto.ExerciseUpdateRequest) (dto.ExerciseResponse, error) {
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

	return utils.ConvertExerciseModelToResponse(ejercicio), nil
}

func (service *ExerciseService) DeleteExercise(id string) error {
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

	return nil
}

func (service *ExerciseService) SearchExercises(search dto.ExerciseSearchRequest) ([]dto.ExerciseResponse, error) {
	ejercicios, err := service.repository.BuscarEjercicios(search.Name, search.Category, search.MuscleGroup, search.Difficulty, search.CreatedBy, search.IncludeDel)
	if err != nil {
		return nil, err
	}

	respuestas := make([]dto.ExerciseResponse, 0, len(ejercicios))
	for _, ejercicio := range ejercicios {
		respuestas = append(respuestas, utils.ConvertExerciseModelToResponse(ejercicio))
	}

	return respuestas, nil
}
