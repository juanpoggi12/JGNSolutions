package repositories

import (
	"context"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoutineRepositoryInterface interface {
	SearchRoutines(name string, userID string, isTemplate *bool, includeDeleted bool) ([]models.Routine, error)
	GetRoutineByID(id string) (models.Routine, error)
	InsertRoutine(routine models.Routine) (*mongo.InsertOneResult, error)
	UpdateRoutine(routine models.Routine) (*mongo.UpdateResult, error)
	DeleteRoutine(id primitive.ObjectID) (*mongo.UpdateResult, error)
}

type RoutineRepository struct {
	db *mongo.Database
}

func NewRoutineRepository(db *mongo.Database) *RoutineRepository {
	return &RoutineRepository{db: db}
}

func (repository RoutineRepository) collection() *mongo.Collection {
	return repository.db.Collection("routines")
}

func (repository RoutineRepository) SearchRoutines(name string, userID string, isTemplate *bool, includeDeleted bool) ([]models.Routine, error) {
	collection := repository.collection()

	filtro := bson.M{}
	if name != "" {
		filtro["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if userID != "" {
		objectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return nil, err
		}
		filtro["userId"] = objectID
	}
	if isTemplate != nil {
		filtro["isTemplate"] = *isTemplate
	}
	if !includeDeleted {
		filtro["isDeleted"] = false
	}

	cursor, err := collection.Find(context.TODO(), filtro)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var routines []models.Routine
	for cursor.Next(context.Background()) {
		var routine models.Routine
		if err := cursor.Decode(&routine); err != nil {
			continue
		}
		routines = append(routines, routine)
	}

	return routines, nil
}

func (repository RoutineRepository) GetRoutineByID(id string) (models.Routine, error) {
	collection := repository.collection()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Routine{}, err
	}

	filtro := bson.M{"_id": objectID}
	var routine models.Routine

	err = collection.FindOne(context.TODO(), filtro).Decode(&routine)
	return routine, err
}

func (repository RoutineRepository) InsertRoutine(routine models.Routine) (*mongo.InsertOneResult, error) {
	collection := repository.collection()
	resultado, err := collection.InsertOne(context.TODO(), routine)
	return resultado, err
}

func (repository RoutineRepository) UpdateRoutine(rutina models.Routine) (*mongo.UpdateResult, error) {
	collection := repository.collection()

	filtro := bson.M{"_id": rutina.ID}
	actualizacion := bson.M{"$set": bson.M{
		"userId":      rutina.UserID,
		"name":        rutina.Name,
		"description": rutina.Description,
		"isTemplate":  rutina.IsTemplate,
		"updatedAt":   rutina.UpdatedAt,
	}}

	resultado, err := collection.UpdateOne(context.TODO(), filtro, actualizacion)
	return resultado, err
}

func (repository RoutineRepository) DeleteRoutine(id primitive.ObjectID) (*mongo.UpdateResult, error) {
	collection := repository.collection()

	filtro := bson.M{"_id": id}
	actualizacion := bson.M{"$set": bson.M{
		"isDeleted": true,
		"updatedAt": time.Now(),
	}}

	resultado, err := collection.UpdateOne(context.TODO(), filtro, actualizacion)
	return resultado, err
}
