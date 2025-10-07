package repositories

import (
	"context"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExerciseRepositoryInterface interface {
	BuscarEjercicios(nombre, categoria, grupoMuscular, dificultad, creadoPor string, incluirEliminados bool) ([]models.Exercise, error)
	ObtenerEjercicioPorID(id string) (models.Exercise, error)
	InsertarEjercicio(ejercicio models.Exercise) (*mongo.InsertOneResult, error)
	ModificarEjercicio(ejercicio models.Exercise) (*mongo.UpdateResult, error)
	EliminarEjercicio(id primitive.ObjectID) (*mongo.UpdateResult, error)
}

type ExerciseRepository struct {
	db *mongo.Database
}

func NewExerciseRepository(db *mongo.Database) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

func (repository ExerciseRepository) collection() *mongo.Collection {
	return repository.db.Collection("exercises")
}

func (repository ExerciseRepository) BuscarEjercicios(nombre, categoria, grupoMuscular, dificultad, creadoPor string, incluirEliminados bool) ([]models.Exercise, error) {
	collection := repository.collection()

	filtro := bson.M{}
	if nombre != "" {
		filtro["name"] = bson.M{"$regex": nombre, "$options": "i"}
	}
	if categoria != "" {
		filtro["category"] = categoria
	}
	if grupoMuscular != "" {
		filtro["muscleGroup"] = grupoMuscular
	}
	if dificultad != "" {
		filtro["difficulty"] = dificultad
	}
	if creadoPor != "" {
		objectID, err := primitive.ObjectIDFromHex(creadoPor)
		if err != nil {
			return nil, err
		}
		filtro["createdByUserId"] = objectID
	}
	if !incluirEliminados {
		filtro["isDeleted"] = false
	}

	cursor, err := collection.Find(context.TODO(), filtro)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var ejercicios []models.Exercise
	for cursor.Next(context.Background()) {
		var ejercicio models.Exercise
		if err := cursor.Decode(&ejercicio); err != nil {
			continue
		}
		ejercicios = append(ejercicios, ejercicio)
	}

	return ejercicios, nil
}

func (repository ExerciseRepository) ObtenerEjercicioPorID(id string) (models.Exercise, error) {
	collection := repository.collection()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Exercise{}, err
	}

	filtro := bson.M{"_id": objectID}
	var ejercicio models.Exercise

	err = collection.FindOne(context.TODO(), filtro).Decode(&ejercicio)
	return ejercicio, err
}

func (repository ExerciseRepository) InsertarEjercicio(ejercicio models.Exercise) (*mongo.InsertOneResult, error) {
	collection := repository.collection()
	resultado, err := collection.InsertOne(context.TODO(), ejercicio)
	return resultado, err
}

func (repository ExerciseRepository) ModificarEjercicio(ejercicio models.Exercise) (*mongo.UpdateResult, error) {
	collection := repository.collection()

	filtro := bson.M{"_id": ejercicio.ID}
	actualizacion := bson.M{"$set": bson.M{
		"name":            ejercicio.Name,
		"description":     ejercicio.Description,
		"category":        ejercicio.Category,
		"muscleGroup":     ejercicio.MuscleGroup,
		"difficulty":      ejercicio.Difficulty,
		"mediaUrl":        ejercicio.MediaURL,
		"instructions":    ejercicio.Instructions,
		"createdByUserId": ejercicio.CreatedBy,
		"updatedAt":       ejercicio.UpdatedAt,
	}}

	resultado, err := collection.UpdateOne(context.TODO(), filtro, actualizacion)
	return resultado, err
}

func (repository ExerciseRepository) EliminarEjercicio(id primitive.ObjectID) (*mongo.UpdateResult, error) {
	collection := repository.collection()

	filtro := bson.M{"_id": id}
	actualizacion := bson.M{"$set": bson.M{
		"isDeleted": true,
		"updatedAt": time.Now(),
	}}

	resultado, err := collection.UpdateOne(context.TODO(), filtro, actualizacion)
	return resultado, err
}
