package repositories

import (
	"context"

	"github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoutineExerciseRepositoryInterface interface {
	BuscarRutinaEjercicios(routineID, exerciseID string) ([]models.RoutineExercise, error)
	ObtenerRutinaEjercicioPorID(id string) (models.RoutineExercise, error)
	InsertarRutinaEjercicio(rutinaEjercicio models.RoutineExercise) (*mongo.InsertOneResult, error)
	ModificarRutinaEjercicio(rutinaEjercicio models.RoutineExercise) (*mongo.UpdateResult, error)
	EliminarRutinaEjercicio(id primitive.ObjectID) (*mongo.DeleteResult, error)
	ObtenerRutinaEjerciciosPorRutinaID(routineID primitive.ObjectID) ([]models.RoutineExercise, error)
}

type RoutineExerciseRepository struct {
	db *mongo.Database
}

func NewRoutineExerciseRepository(db *mongo.Database) *RoutineExerciseRepository {
	return &RoutineExerciseRepository{db: db}
}

func (repository RoutineExerciseRepository) collection() *mongo.Collection {
	return repository.db.Collection("routine_exercises")
}

func (repository RoutineExerciseRepository) BuscarRutinaEjercicios(routineID, exerciseID string) ([]models.RoutineExercise, error) {
	collection := repository.collection()

	filtro := bson.M{}
	if routineID != "" {
		objectID, err := primitive.ObjectIDFromHex(routineID)
		if err != nil {
			return nil, err
		}
		filtro["routineId"] = objectID
	}
	if exerciseID != "" {
		objectID, err := primitive.ObjectIDFromHex(exerciseID)
		if err != nil {
			return nil, err
		}
		filtro["exerciseId"] = objectID
	}

	cursor, err := collection.Find(context.TODO(), filtro)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var rutinasEjercicio []models.RoutineExercise
	for cursor.Next(context.Background()) {
		var rutina models.RoutineExercise
		if err := cursor.Decode(&rutina); err != nil {
			continue
		}
		rutinasEjercicio = append(rutinasEjercicio, rutina)
	}

	return rutinasEjercicio, nil
}

func (repository RoutineExerciseRepository) ObtenerRutinaEjercicioPorID(id string) (models.RoutineExercise, error) {
	collection := repository.collection()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.RoutineExercise{}, err
	}

	filtro := bson.M{"_id": objectID}
	var rutina models.RoutineExercise

	err = collection.FindOne(context.TODO(), filtro).Decode(&rutina)
	return rutina, err
}

func (repository RoutineExerciseRepository) InsertarRutinaEjercicio(rutinaEjercicio models.RoutineExercise) (*mongo.InsertOneResult, error) {
	collection := repository.collection()
	resultado, err := collection.InsertOne(context.TODO(), rutinaEjercicio)
	return resultado, err
}

func (repository RoutineExerciseRepository) ModificarRutinaEjercicio(rutinaEjercicio models.RoutineExercise) (*mongo.UpdateResult, error) {
	collection := repository.collection()

	filtro := bson.M{"_id": rutinaEjercicio.ID}
	actualizacion := bson.M{"$set": bson.M{
		"routineId":     rutinaEjercicio.RoutineID,
		"exerciseId":    rutinaEjercicio.ExerciseID,
		"order":         rutinaEjercicio.Order,
		"sets":          rutinaEjercicio.Sets,
		"reps":          rutinaEjercicio.Reps,
		"targetWeight":  rutinaEjercicio.TargetWeight,
		"targetTimeSec": rutinaEjercicio.TargetTimeSec,
		"notes":         rutinaEjercicio.Notes,
	}}

	resultado, err := collection.UpdateOne(context.TODO(), filtro, actualizacion)
	return resultado, err
}

func (repository RoutineExerciseRepository) EliminarRutinaEjercicio(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	collection := repository.collection()
	resultado, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return resultado, err
}

func (repository RoutineExerciseRepository) ObtenerRutinaEjerciciosPorRutinaID(routineID primitive.ObjectID) ([]models.RoutineExercise, error) {
	collection := repository.collection()

	cursor, err := collection.Find(context.TODO(), bson.M{"routineId": routineID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var rutinasEjercicio []models.RoutineExercise
	for cursor.Next(context.Background()) {
		var rutina models.RoutineExercise
		if err := cursor.Decode(&rutina); err != nil {
			continue
		}
		rutinasEjercicio = append(rutinasEjercicio, rutina)
	}

	return rutinasEjercicio, nil
}
