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
	GetRoutineByID(actor Actor, id string) (dto.RoutineResponse, error)
	UpdateRoutine(actor Actor, id string, req dto.RoutineUpdateRequest) (dto.RoutineResponse, error)
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
}

func NewRoutineService(r repositories.RoutineRepositoryInterface, re repositories.RoutineExerciseRepositoryInterface, e repositories.ExerciseRepositoryInterface) *RoutineService {
	return &RoutineService{
		routineRepository:         r,
		routineExerciseRepository: re,
		exerciseRepository:        e,
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

func (service *RoutineService) GetRoutineByID(actor Actor, id string) (dto.RoutineResponse, error) {
	rutina, err := service.routineRepository.ObtenerRutinaPorID(id)
	if err != nil {
		return dto.RoutineResponse{}, errors.New("rutina no encontrada")
	}
	if rutina.IsDeleted {
		return dto.RoutineResponse{}, errors.New("la rutina fue eliminada")
	}

	// Solo el dueño o un admin puede acceder
	if actor.Role != "admin" && rutina.UserID != actor.UserID {
		return dto.RoutineResponse{}, errors.New("no tienes permiso para acceder a esta rutina")
	}

	return utils.ConvertRoutineModelToResponse(rutina), nil
}

func (service *RoutineService) UpdateRoutine(actor Actor, id string, req dto.RoutineUpdateRequest) (dto.RoutineResponse, error) {
	rutina, err := service.routineRepository.ObtenerRutinaPorID(id)
	if err != nil {
		return dto.RoutineResponse{}, errors.New("rutina no encontrada")
	}
	if rutina.IsDeleted {
		return dto.RoutineResponse{}, errors.New("la rutina fue eliminada")
	}

	if actor.Role != "admin" && rutina.UserID != actor.UserID {
		return dto.RoutineResponse{}, errors.New("no tienes permiso para modificar esta rutina")
	}

	utils.ApplyRoutineUpdateToModel(&rutina, req)
	rutina.UpdatedAt = time.Now()

	_, err = service.routineRepository.ModificarRutina(rutina)
	if err != nil {
		return dto.RoutineResponse{}, err
	}

	return utils.ConvertRoutineModelToResponse(rutina), nil
}

func (service *RoutineService) DeleteRoutine(actor Actor, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID inválido")
	}

	rutina, err := service.routineRepository.ObtenerRutinaPorID(id)
	if err != nil {
		return errors.New("rutina no encontrada")
	}

	if actor.Role != "admin" && rutina.UserID != actor.UserID {
		return errors.New("no tienes permiso para eliminar esta rutina")
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
func (service *RoutineService) SearchRoutines(actor Actor, search dto.RoutineSearchRequest) ([]dto.RoutineResponse, error) {
	var userID string

	// Un usuario solo puede ver sus propias rutinas
	if actor.Role != "admin" {
		userID = actor.UserID.Hex()
	}

	rutinas, err := service.routineRepository.BuscarRutinas(
		search.Name,
		userID, // usamos el userID derivado del actor, no del DTO
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

func (service *RoutineService) UpdateRoutineExercise(actor Actor, id string, req dto.RoutineExerciseUpdateRequest) (dto.RoutineExerciseResponse, error) {
	rutinaEjercicio, err := service.routineExerciseRepository.ObtenerRutinaEjercicioPorID(id)
	if err != nil {
		return dto.RoutineExerciseResponse{}, errors.New("asociación rutina-ejercicio no encontrada")
	}

	rutina, err := service.routineRepository.ObtenerRutinaPorID(rutinaEjercicio.RoutineID.Hex())
	if err != nil {
		return dto.RoutineExerciseResponse{}, errors.New("rutina no encontrada")
	}

	if actor.Role != "admin" && rutina.UserID != actor.UserID {
		return dto.RoutineExerciseResponse{}, errors.New("no tienes permiso para modificar ejercicios en esta rutina")
	}

	utils.ApplyRoutineExerciseUpdateToModel(&rutinaEjercicio, req)

	_, err = service.routineExerciseRepository.ModificarRutinaEjercicio(rutinaEjercicio)
	if err != nil {
		return dto.RoutineExerciseResponse{}, err
	}

	return utils.ConvertRoutineExerciseModelToResponse(rutinaEjercicio), nil
}

func (service *RoutineService) DeleteRoutineExercise(actor Actor, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID inválido")
	}

	rutinaEjercicio, err := service.routineExerciseRepository.ObtenerRutinaEjercicioPorID(id)
	if err != nil {
		return errors.New("asociación rutina-ejercicio no encontrada")
	}

	rutina, err := service.routineRepository.ObtenerRutinaPorID(rutinaEjercicio.RoutineID.Hex())
	if err != nil {
		return errors.New("rutina no encontrada")
	}

	if actor.Role != "admin" && rutina.UserID != actor.UserID {
		return errors.New("no tienes permiso para eliminar ejercicios de esta rutina")
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

func (service *RoutineService) ListRoutineExercises(actor Actor, search dto.RoutineExerciseSearchRequest) ([]dto.RoutineExerciseResponse, error) {
	rutina, err := service.cargarRutinaActiva(search.RoutineID)
	if err != nil {
		return nil, err
	}

	if actor.Role != "admin" && rutina.UserID != actor.UserID {
		return nil, errors.New("no tienes permiso para ver los ejercicios de esta rutina")
	}

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

func (service *RoutineService) DuplicateRoutine(actor Actor, routineID string, newName string) (dto.RoutineResponse, error) {
	original, err := service.cargarRutinaActiva(routineID)
	if err != nil {
		return dto.RoutineResponse{}, err
	}

	// Solo admin o el dueño puede duplicar la rutina
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

// --- Funciones privadas (no cambian) ---
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
