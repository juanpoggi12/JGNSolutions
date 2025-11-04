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
	CreateRoutine(actor Actor, req dto.RoutineCreateRequest) (dto.RoutineResponse, error)
	DeleteRoutine(actor Actor, id string) error
	SearchRoutines(actor Actor, search dto.RoutineSearchRequest) ([]dto.RoutineResponse, error)
	AddExerciseToRoutine(actor Actor, req dto.RoutineExerciseCreateRequest) (dto.RoutineExerciseResponse, error)
	UpdateRoutineExercise(actor Actor, id string, req dto.RoutineExerciseUpdateRequest) (dto.RoutineExerciseResponse, error)
	DeleteRoutineExercise(actor Actor, id string) error
	ListRoutineExercises(actor Actor, search dto.RoutineExerciseSearchRequest) ([]dto.RoutineExerciseResponse, error)
	DuplicateRoutine(actor Actor, id string, newName string) (dto.RoutineResponse, error)
}

type RoutineService struct {
	routineRepository         repositories.RoutineRepositoryInterface
	routineExerciseRepository repositories.RoutineExerciseRepositoryInterface
	exerciseRepository        repositories.ExerciseRepositoryInterface
	logService                *LogService
}

func NewRoutineService(r repositories.RoutineRepositoryInterface, re repositories.RoutineExerciseRepositoryInterface, e repositories.ExerciseRepositoryInterface, logService *LogService) *RoutineService {
	return &RoutineService{
		routineRepository:         r,
		routineExerciseRepository: re,
		exerciseRepository:        e,
		logService:                logService,
	}
}

func (service *RoutineService) CreateRoutine(actor Actor, req dto.RoutineCreateRequest) (dto.RoutineResponse, error) {
	rutina, err := utils.ConvertRoutineCreateRequestToModel(req)
	if err != nil {
		return dto.RoutineResponse{}, err
	}

	// Asignar propietario según el actor autenticado
	rutina.UserID = actor.UserID
	rutina.CreatedAt = time.Now()
	rutina.UpdatedAt = time.Now()
	rutina.IsDeleted = false

	resultado, err := service.routineRepository.InsertRoutine(rutina)
	if err != nil {
		return dto.RoutineResponse{}, err
	}

	if oid, ok := resultado.InsertedID.(primitive.ObjectID); ok {
		rutina.ID = oid

		// Log de creación exitoso
		if service.logService != nil {
			service.logService.RecordAction(actor.UserID, "rutina:creada "+rutina.Name)
		}

		return utils.ConvertRoutineModelToResponse(rutina), nil
	}

	return dto.RoutineResponse{}, errors.New("error al obtener ID de la rutina creada")
}

func (service *RoutineService) DeleteRoutine(actor Actor, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID inválido")
	}

	rutina, err := service.routineRepository.GetRoutineByID(id)
	if err != nil {
		return errors.New("rutina no encontrada")
	}

	if actor.Role != "admin" && rutina.UserID != actor.UserID {
		return errors.New("no tienes permiso para eliminar esta rutina")
	}

	resultado, err := service.routineRepository.DeleteRoutine(objectID)
	if err != nil {
		return err
	}
	if resultado.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	// Log de eliminación exitosa
	if service.logService != nil {
		service.logService.RecordAction(actor.UserID, "rutina:eliminada "+id)
	}

	return nil
}

func (service *RoutineService) SearchRoutines(actor Actor, search dto.RoutineSearchRequest) ([]dto.RoutineResponse, error) {
	var userID string

	if actor.Role != "admin" {
		userID = actor.UserID.Hex()
	}

	rutinas, err := service.routineRepository.SearchRoutines(
		search.Name,
		userID,
		search.IsTemplate,
		search.IncludeDel,
	)
	if err != nil {
		return nil, err
	}

	respuestas := make([]dto.RoutineResponse, 0, len(rutinas))
	for _, rutina := range rutinas {
		respuestas = append(respuestas, utils.ConvertRoutineModelToResponse(rutina))
	}

	return respuestas, nil
}

func (service *RoutineService) AddExerciseToRoutine(actor Actor, req dto.RoutineExerciseCreateRequest) (dto.RoutineExerciseResponse, error) {
	rutina, err := service.cargarRutinaActiva(req.RoutineID)
	if err != nil {
		return dto.RoutineExerciseResponse{}, err
	}

	if actor.Role != "admin" && rutina.UserID != actor.UserID {
		return dto.RoutineExerciseResponse{}, errors.New("no tienes permiso para modificar esta rutina")
	}

	if _, err := service.cargarEjercicioActivo(req.ExerciseID); err != nil {
		return dto.RoutineExerciseResponse{}, err
	}

	rutinaEjercicio, err := utils.ConvertRoutineExerciseCreateRequestToModel(req)
	if err != nil {
		return dto.RoutineExerciseResponse{}, err
	}

	resultado, err := service.routineExerciseRepository.InsertRoutineExercise(rutinaEjercicio)
	if err != nil {
		return dto.RoutineExerciseResponse{}, err
	}

	if oid, ok := resultado.InsertedID.(primitive.ObjectID); ok {
		rutinaEjercicio.ID = oid

		if service.logService != nil {
			service.logService.RecordAction(actor.UserID, "rutina_ejercicio:agregado "+rutinaEjercicio.RoutineID.Hex())
		}

		return utils.ConvertRoutineExerciseModelToResponse(rutinaEjercicio), nil
	}

	return dto.RoutineExerciseResponse{}, errors.New("error al crear rutina ejercicio")
}

func (service *RoutineService) UpdateRoutineExercise(actor Actor, id string, req dto.RoutineExerciseUpdateRequest) (dto.RoutineExerciseResponse, error) {
	rutinaEjercicio, err := service.routineExerciseRepository.GetRoutineExerciseByID(id)
	if err != nil {
		return dto.RoutineExerciseResponse{}, errors.New("asociación rutina-ejercicio no encontrada")
	}

	rutina, err := service.routineRepository.GetRoutineByID(rutinaEjercicio.RoutineID.Hex())
	if err != nil {
		return dto.RoutineExerciseResponse{}, errors.New("rutina no encontrada")
	}

	if actor.Role != "admin" && rutina.UserID != actor.UserID {
		return dto.RoutineExerciseResponse{}, errors.New("no tienes permiso para modificar ejercicios en esta rutina")
	}

	utils.ApplyRoutineExerciseUpdateToModel(&rutinaEjercicio, req)

	_, err = service.routineExerciseRepository.UpdateRoutineExercise(rutinaEjercicio)
	if err != nil {
		return dto.RoutineExerciseResponse{}, err
	}

	if service.logService != nil {
		service.logService.RecordAction(actor.UserID, "rutina_ejercicio:modificado "+rutinaEjercicio.ID.Hex())
	}

	return utils.ConvertRoutineExerciseModelToResponse(rutinaEjercicio), nil
}

func (service *RoutineService) DeleteRoutineExercise(actor Actor, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID inválido")
	}

	rutinaEjercicio, err := service.routineExerciseRepository.GetRoutineExerciseByID(id)
	if err != nil {
		return errors.New("asociación rutina-ejercicio no encontrada")
	}

	rutina, err := service.routineRepository.GetRoutineByID(rutinaEjercicio.RoutineID.Hex())
	if err != nil {
		return errors.New("rutina no encontrada")
	}

	if actor.Role != "admin" && rutina.UserID != actor.UserID {
		return errors.New("no tienes permiso para eliminar ejercicios de esta rutina")
	}

	resultado, err := service.routineExerciseRepository.DeleteRoutineExercise(objectID)
	if err != nil {
		return err
	}
	if resultado.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	if service.logService != nil {
		service.logService.RecordAction(actor.UserID, "rutina_ejercicio:eliminado "+id)
	}

	return nil
}

func (service *RoutineService) ListRoutineExercises(actor Actor, search dto.RoutineExerciseSearchRequest) ([]dto.RoutineExerciseResponse, error) {
	rutina, err := service.cargarRutinaActiva(search.RoutineID)
	if err != nil {
		return nil, err
	}

	if actor.Role != "admin" && rutina.UserID != actor.UserID {
		return nil, errors.New("no tienes permiso para ver los ejercicios de esta rutina")
	}

	rutinasEjercicio, err := service.routineExerciseRepository.SearchRoutineExercises(search.RoutineID, search.ExerciseID)
	if err != nil {
		return nil, err
	}

	respuestas := make([]dto.RoutineExerciseResponse, 0, len(rutinasEjercicio))
	for _, rutina := range rutinasEjercicio {
		respuestas = append(respuestas, utils.ConvertRoutineExerciseModelToResponse(rutina))
	}

	return respuestas, nil
}

func (service *RoutineService) DuplicateRoutine(actor Actor, routineID string, newName string) (dto.RoutineResponse, error) {
	original, err := service.cargarRutinaActiva(routineID)
	if err != nil {
		return dto.RoutineResponse{}, err
	}

	if actor.Role != "admin" && original.UserID != actor.UserID {
		return dto.RoutineResponse{}, errors.New("no tienes permiso para duplicar esta rutina")
	}

	nuevaRutina := models.Routine{
		UserID:      actor.UserID,
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

	resultado, err := service.routineRepository.InsertRoutine(nuevaRutina)
	if err != nil {
		return dto.RoutineResponse{}, err
	}

	if oid, ok := resultado.InsertedID.(primitive.ObjectID); ok {
		nuevaRutina.ID = oid
	} else {
		return dto.RoutineResponse{}, errors.New("error al duplicar rutina")
	}

	ejerciciosOriginales, err := service.routineExerciseRepository.GetRoutineExercisesByRoutineID(original.ID)
	if err != nil {
		return dto.RoutineResponse{}, err
	}

	for _, rutinaEjercicio := range ejerciciosOriginales {
		rutinaEjercicio.ID = primitive.NilObjectID
		rutinaEjercicio.RoutineID = nuevaRutina.ID
		if _, err := service.routineExerciseRepository.InsertRoutineExercise(rutinaEjercicio); err != nil {
			return dto.RoutineResponse{}, err
		}
	}

	if service.logService != nil {
		service.logService.RecordAction(actor.UserID, "rutina:duplicada "+nuevaRutina.Name)
	}

	return utils.ConvertRoutineModelToResponse(nuevaRutina), nil
}

func (service *RoutineService) cargarRutinaActiva(id string) (models.Routine, error) {
	rutina, err := service.routineRepository.GetRoutineByID(id)
	if err != nil {
		return models.Routine{}, errors.New("rutina no encontrada")
	}
	if rutina.IsDeleted {
		return models.Routine{}, errors.New("la rutina fue eliminada")
	}
	return rutina, nil
}

func (service *RoutineService) cargarEjercicioActivo(id string) (models.Exercise, error) {
	ejercicio, err := service.exerciseRepository.GetExerciseByID(id)
	if err != nil {
		return models.Exercise{}, errors.New("ejercicio no encontrado")
	}
	if ejercicio.IsDeleted {
		return models.Exercise{}, errors.New("el ejercicio fue eliminado")
	}
	return ejercicio, nil
}
