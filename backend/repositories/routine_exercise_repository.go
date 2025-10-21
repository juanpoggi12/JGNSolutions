package repositories

import (
	"context"

	"github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoutineExerciseRepositoryInterface interface {
	SearchRoutineExercises(routineID, exerciseID string) ([]models.RoutineExercise, error)
	GetRoutineExerciseByID(id string) (models.RoutineExercise, error)
	InsertRoutineExercise(re models.RoutineExercise) (*mongo.InsertOneResult, error)
	UpdateRoutineExercise(re models.RoutineExercise) (*mongo.UpdateResult, error)
	DeleteRoutineExercise(id primitive.ObjectID) (*mongo.DeleteResult, error)
	GetRoutineExercisesByRoutineID(routineID primitive.ObjectID) ([]models.RoutineExercise, error)
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

func (repository RoutineExerciseRepository) SearchRoutineExercises(routineID, exerciseID string) ([]models.RoutineExercise, error) {
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

	var routineExercises []models.RoutineExercise
	for cursor.Next(context.Background()) {
		var re models.RoutineExercise
		if err := cursor.Decode(&re); err != nil {
			continue
		}
		routineExercises = append(routineExercises, re)
	}

	return routineExercises, nil
}

func (repository RoutineExerciseRepository) GetRoutineExerciseByID(id string) (models.RoutineExercise, error) {
	collection := repository.collection()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.RoutineExercise{}, err
	}

	filtro := bson.M{"_id": objectID}
	var re models.RoutineExercise

	err = collection.FindOne(context.TODO(), filtro).Decode(&re)
	return re, err
}

func (repository RoutineExerciseRepository) InsertRoutineExercise(re models.RoutineExercise) (*mongo.InsertOneResult, error) {
	collection := repository.collection()
	resultado, err := collection.InsertOne(context.TODO(), re)
	return resultado, err
}

func (repository RoutineExerciseRepository) UpdateRoutineExercise(re models.RoutineExercise) (*mongo.UpdateResult, error) {
	collection := repository.collection()

	filtro := bson.M{"_id": re.ID}
	actualizacion := bson.M{"$set": bson.M{
		"routineId":     re.RoutineID,
		"exerciseId":    re.ExerciseID,
		"order":         re.Order,
		"sets":          re.Sets,
		"reps":          re.Reps,
		"targetWeight":  re.TargetWeight,
		"targetTimeSec": re.TargetTimeSec,
		"notes":         re.Notes,
	}}

	resultado, err := collection.UpdateOne(context.TODO(), filtro, actualizacion)
	return resultado, err
}

func (repository RoutineExerciseRepository) DeleteRoutineExercise(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	collection := repository.collection()
	resultado, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return resultado, err
}

func (repository RoutineExerciseRepository) GetRoutineExercisesByRoutineID(routineID primitive.ObjectID) ([]models.RoutineExercise, error) {
	collection := repository.collection()

	cursor, err := collection.Find(context.TODO(), bson.M{"routineId": routineID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var routineExercises []models.RoutineExercise
	for cursor.Next(context.Background()) {
		var re models.RoutineExercise
		if err := cursor.Decode(&re); err != nil {
			continue
		}
		routineExercises = append(routineExercises, re)
	}

	return routineExercises, nil
}
