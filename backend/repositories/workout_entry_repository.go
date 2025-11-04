package repositories

import (
	"context"

	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Interfaz del repositorio (contrato que usa el service)
type WorkoutEntryRepositoryInterface interface {
	Create(entry *models.WorkoutEntry) error
	Search(filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutEntry, error)
}

// Implementación concreta (usa MongoDB)
type WorkoutEntryRepository struct {
	collection *mongo.Collection
}

// Constructor
func NewWorkoutEntryRepository(db *mongo.Database) *WorkoutEntryRepository {
	return &WorkoutEntryRepository{collection: db.Collection("workoutEntries")}
}

// Crear una entrada
func (r *WorkoutEntryRepository) Create(entry *models.WorkoutEntry) error {
	res, err := r.collection.InsertOne(context.TODO(), entry)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		entry.ID = oid
	}
	return nil
}

// Buscar con filtros
func (r *WorkoutEntryRepository) Search(filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutEntry, error) {
	cursor, err := r.collection.Find(context.TODO(), filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var entries []models.WorkoutEntry
	if err := cursor.All(context.TODO(), &entries); err != nil {
		return nil, err
	}
	return entries, nil
}
