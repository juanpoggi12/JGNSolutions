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
	BuscarRutinas(nombre string, userID string, esPlantilla *bool, incluirEliminadas bool) ([]models.Routine, error)
	ObtenerRutinaPorID(id string) (models.Routine, error)
	InsertarRutina(rutina models.Routine) (*mongo.InsertOneResult, error)
	ModificarRutina(rutina models.Routine) (*mongo.UpdateResult, error)
	EliminarRutina(id primitive.ObjectID) (*mongo.UpdateResult, error)
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

func (repository RoutineRepository) BuscarRutinas(nombre string, userID string, esPlantilla *bool, incluirEliminadas bool) ([]models.Routine, error) {
	collection := repository.collection()

	filtro := bson.M{}
	if nombre != "" {
		filtro["name"] = bson.M{"$regex": nombre, "$options": "i"}
	}
	if userID != "" {
		objectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return nil, err
		}
		filtro["userId"] = objectID
	}
	if esPlantilla != nil {
		filtro["isTemplate"] = *esPlantilla
	}
	if !incluirEliminadas {
		filtro["isDeleted"] = false
	}

	cursor, err := collection.Find(context.TODO(), filtro)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var rutinas []models.Routine
	for cursor.Next(context.Background()) {
		var rutina models.Routine
		if err := cursor.Decode(&rutina); err != nil {
			continue
		}
		rutinas = append(rutinas, rutina)
	}

	return rutinas, nil
}

func (repository RoutineRepository) ObtenerRutinaPorID(id string) (models.Routine, error) {
	collection := repository.collection()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Routine{}, err
	}

	filtro := bson.M{"_id": objectID}
	var rutina models.Routine

	err = collection.FindOne(context.TODO(), filtro).Decode(&rutina)
	return rutina, err
}

func (repository RoutineRepository) InsertarRutina(rutina models.Routine) (*mongo.InsertOneResult, error) {
	collection := repository.collection()
	resultado, err := collection.InsertOne(context.TODO(), rutina)
	return resultado, err
}

func (repository RoutineRepository) ModificarRutina(rutina models.Routine) (*mongo.UpdateResult, error) {
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

func (repository RoutineRepository) EliminarRutina(id primitive.ObjectID) (*mongo.UpdateResult, error) {
	collection := repository.collection()

	filtro := bson.M{"_id": id}
	actualizacion := bson.M{"$set": bson.M{
		"isDeleted": true,
		"updatedAt": time.Now(),
	}}

	resultado, err := collection.UpdateOne(context.TODO(), filtro, actualizacion)
	return resultado, err
}
