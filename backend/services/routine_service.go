package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/repositories"
	"github.com/juanpoggi12/JGNSolutions/backend/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoutineServiceInterface interface {
	CreateRoutine(req dto.RoutineCreateRequest) (dto.RoutineResponse, error)
	GetRoutineByID(id string) (dto.RoutineResponse, error)
	UpdateRoutine(id string, req dto.RoutineUpdateRequest) (dto.RoutineResponse, error)
	DeleteRoutine(id string) error
	SearchRoutines(search dto.RoutineSearchRequest) ([]dto.RoutineResponse, error)
	AddExerciseToRoutine(req dto.RoutineExerciseCreateRequest) (dto.RoutineExerciseResponse, error)
	UpdateRoutineExercise(id string, req dto.RoutineExerciseUpdateRequest) (dto.RoutineExerciseResponse, error)
	DeleteRoutineExercise(id string) error
	ListRoutineExercises(search dto.RoutineExerciseSearchRequest) ([]dto.RoutineExerciseResponse, error)
	DuplicateRoutine(id string, newName string, targetUserID string) (dto.RoutineResponse, error)
}

type RoutineService struct {
	routineRepository         repositories.RoutineRepositoryInterface
	routineExerciseRepository repositories.RoutineExerciseRepositoryInterface
	exerciseRepository        repositories.ExerciseRepositoryInterface
}

func NewRoutineService(r repositories.RoutineRepositoryInterface, re repositories.RoutineExerciseRepositoryInterface, e repositories.ExerciseRepositoryInterface) *RoutineService {
	return &RoutineService{
		routineRepository:         r,
		routineExerciseRepository: re,
		exerciseRepository:        e,
	}
}

func (service *RoutineService) CreateRoutine(req dto.RoutineCreateRequest) (dto.RoutineResponse, error) {
	rutina, err := utils.ConvertRoutineCreateRequestToModel(req)
	if err != nil {
		return dto.RoutineResponse{}, err
	}

	resultado, err := service.routineRepository.InsertarRutina(rutina)
	if err != nil {
		return dto.RoutineResponse{}, err
	}

	if oid, ok := resultado.InsertedID.(primitive.ObjectID); ok {
		rutina.ID = oid
		return utils.ConvertRoutineModelToResponse(rutina), nil
	}

	return dto.RoutineResponse{}, errors.New("error al obtener ID de la rutina creada")
}

func (service *RoutineService) GetRoutineByID(id string) (dto.RoutineResponse, error) {
	rutina, err := service.routineRepository.ObtenerRutinaPorID(id)
	if err != nil {
		return dto.RoutineResponse{}, errors.New("rutina no encontrada")
	}
	if rutina.IsDeleted {
		return dto.RoutineResponse{}, errors.New("la rutina fue eliminada")
	}

	return utils.ConvertRoutineModelToResponse(rutina), nil
}

func (service *RoutineService) UpdateRoutine(id string, req dto.RoutineUpdateRequest) (dto.RoutineResponse, error) {
	rutina, err := service.routineRepository.ObtenerRutinaPorID(id)
	if err != nil {
		return dto.RoutineResponse{}, errors.New("rutina no encontrada")
	}
	if rutina.IsDeleted {
		return dto.RoutineResponse{}, errors.New("la rutina fue eliminada")
	}

	utils.ApplyRoutineUpdateToModel(&rutina, req)
	rutina.UpdatedAt = time.Now()

	_, err = service.routineRepository.ModificarRutina(rutina)
	if err != nil {
		return dto.RoutineResponse{}, err
	}

	return utils.ConvertRoutineModelToResponse(rutina), nil
}

func (service *RoutineService) DeleteRoutine(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID inválido")
	}

	resultado, err := service.routineRepository.EliminarRutina(objectID)
	if err != nil {
		return err
	}
	if resultado.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (service *RoutineService) SearchRoutines(search dto.RoutineSearchRequest) ([]dto.RoutineResponse, error) {
	rutinas, err := service.routineRepository.BuscarRutinas(search.Name, search.UserID, search.IsTemplate, search.IncludeDel)
	if err != nil {
		return nil, err
	}

	respuestas := make([]dto.RoutineResponse, 0, len(rutinas))
	for _, rutina := range rutinas {
		respuestas = append(respuestas, utils.ConvertRoutineModelToResponse(rutina))
	}

	return respuestas, nil
}

func (service *RoutineService) AddExerciseToRoutine(req dto.RoutineExerciseCreateRequest) (dto.RoutineExerciseResponse, error) {
	if _, err := service.cargarRutinaActiva(req.RoutineID); err != nil {
		return dto.RoutineExerciseResponse{}, err
	}
	if _, err := service.cargarEjercicioActivo(req.ExerciseID); err != nil {
		return dto.RoutineExerciseResponse{}, err
	}

	rutinaEjercicio, err := utils.ConvertRoutineExerciseCreateRequestToModel(req)
	if err != nil {
		return dto.RoutineExerciseResponse{}, err
	}

	resultado, err := service.routineExerciseRepository.InsertarRutinaEjercicio(rutinaEjercicio)
	if err != nil {
		return dto.RoutineExerciseResponse{}, err
	}

	if oid, ok := resultado.InsertedID.(primitive.ObjectID); ok {
		rutinaEjercicio.ID = oid
		return utils.ConvertRoutineExerciseModelToResponse(rutinaEjercicio), nil
	}

	return dto.RoutineExerciseResponse{}, errors.New("error al crear rutina ejercicio")
}

func (service *RoutineService) UpdateRoutineExercise(id string, req dto.RoutineExerciseUpdateRequest) (dto.RoutineExerciseResponse, error) {
	rutinaEjercicio, err := service.routineExerciseRepository.ObtenerRutinaEjercicioPorID(id)
	if err != nil {
		return dto.RoutineExerciseResponse{}, errors.New("asociación rutina-ejercicio no encontrada")
	}

	utils.ApplyRoutineExerciseUpdateToModel(&rutinaEjercicio, req)

	_, err = service.routineExerciseRepository.ModificarRutinaEjercicio(rutinaEjercicio)
	if err != nil {
		return dto.RoutineExerciseResponse{}, err
	}

	return utils.ConvertRoutineExerciseModelToResponse(rutinaEjercicio), nil
}

func (service *RoutineService) DeleteRoutineExercise(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID inválido")
	}

	resultado, err := service.routineExerciseRepository.EliminarRutinaEjercicio(objectID)
	if err != nil {
		return err
	}
	if resultado.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (service *RoutineService) ListRoutineExercises(search dto.RoutineExerciseSearchRequest) ([]dto.RoutineExerciseResponse, error) {
	rutinasEjercicio, err := service.routineExerciseRepository.BuscarRutinaEjercicios(search.RoutineID, search.ExerciseID)
	if err != nil {
		return nil, err
	}

	respuestas := make([]dto.RoutineExerciseResponse, 0, len(rutinasEjercicio))
	for _, rutina := range rutinasEjercicio {
		respuestas = append(respuestas, utils.ConvertRoutineExerciseModelToResponse(rutina))
	}

	return respuestas, nil
}

func (service *RoutineService) DuplicateRoutine(routineID string, newName string, targetUserID string) (dto.RoutineResponse, error) {
	original, err := service.cargarRutinaActiva(routineID)
	if err != nil {
		return dto.RoutineResponse{}, err
	}

	userID := original.UserID
	if targetUserID != "" {
		oid, err := primitive.ObjectIDFromHex(targetUserID)
		if err != nil {
			return dto.RoutineResponse{}, errors.New("ID de usuario inválido")
		}
		userID = oid
	}

	nuevaRutina := models.Routine{
		UserID:      userID,
		Name:        original.Name,
		Description: original.Description,
		IsTemplate:  original.IsTemplate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsDeleted:   false,
	}

	if newName != "" {
		nuevaRutina.Name = newName
	} else {
		nuevaRutina.Name = fmt.Sprintf("%s (Copia)", original.Name)
	}

	resultado, err := service.routineRepository.InsertarRutina(nuevaRutina)
	if err != nil {
		return dto.RoutineResponse{}, err
	}

	if oid, ok := resultado.InsertedID.(primitive.ObjectID); ok {
		nuevaRutina.ID = oid
	} else {
		return dto.RoutineResponse{}, errors.New("error al duplicar rutina")
	}

	ejerciciosOriginales, err := service.routineExerciseRepository.ObtenerRutinaEjerciciosPorRutinaID(original.ID)
	if err != nil {
		return dto.RoutineResponse{}, err
	}

	for _, rutinaEjercicio := range ejerciciosOriginales {
		rutinaEjercicio.ID = primitive.NilObjectID
		rutinaEjercicio.RoutineID = nuevaRutina.ID
		if _, err := service.routineExerciseRepository.InsertarRutinaEjercicio(rutinaEjercicio); err != nil {
			return dto.RoutineResponse{}, err
		}
	}

	return utils.ConvertRoutineModelToResponse(nuevaRutina), nil
}

func (service *RoutineService) cargarRutinaActiva(id string) (models.Routine, error) {
	rutina, err := service.routineRepository.ObtenerRutinaPorID(id)
	if err != nil {
		return models.Routine{}, errors.New("rutina no encontrada")
	}
	if rutina.IsDeleted {
		return models.Routine{}, errors.New("la rutina fue eliminada")
	}
	return rutina, nil
}

func (service *RoutineService) cargarEjercicioActivo(id string) (models.Exercise, error) {
	ejercicio, err := service.exerciseRepository.ObtenerEjercicioPorID(id)
	if err != nil {
		return models.Exercise{}, errors.New("ejercicio no encontrado")
	}
	if ejercicio.IsDeleted {
		return models.Exercise{}, errors.New("el ejercicio fue eliminado")
	}
	return ejercicio, nil
}
