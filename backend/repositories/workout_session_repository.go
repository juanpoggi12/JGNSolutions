package repositories

import (
	"context"

	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 🧩 Interfaz del repositorio (contrato que usa el service)
type WorkoutSessionRepositoryInterface interface {
	Create(session *models.WorkoutSession) error
	FindByID(id primitive.ObjectID) (*models.WorkoutSession, error)
	FindByUser(userID primitive.ObjectID) ([]models.WorkoutSession, error)
	Update(session *models.WorkoutSession) error
	Delete(id primitive.ObjectID) error
	Search(filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutSession, error)
}

// 🧩 Implementación concreta (usa MongoDB)
type WorkoutSessionRepository struct {
	collection *mongo.Collection
}

// Constructor
func NewWorkoutSessionRepository(db *mongo.Database) *WorkoutSessionRepository {
	return &WorkoutSessionRepository{collection: db.Collection("workoutSessions")}
}

// Crear una sesión
func (r *WorkoutSessionRepository) Create(session *models.WorkoutSession) error {
	res, err := r.collection.InsertOne(context.TODO(), session)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		session.ID = oid
	}

	return nil
}

// Buscar por ID
func (r *WorkoutSessionRepository) FindByID(id primitive.ObjectID) (*models.WorkoutSession, error) {
	var session models.WorkoutSession
	if err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&session); err != nil {
		return nil, err
	}
	return &session, nil
}

// Buscar todas las sesiones de un usuario
func (r *WorkoutSessionRepository) FindByUser(userID primitive.ObjectID) ([]models.WorkoutSession, error) {
	var sessions []models.WorkoutSession

	cursor, err := r.collection.Find(context.TODO(), bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

// Actualizar una sesión
func (r *WorkoutSessionRepository) Update(session *models.WorkoutSession) error {
	_, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": session.ID}, bson.M{"$set": session})
	return err
}

// Eliminar una sesión
func (r *WorkoutSessionRepository) Delete(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

// Buscar con filtros personalizados
func (r *WorkoutSessionRepository) Search(filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutSession, error) {
	cursor, err := r.collection.Find(context.TODO(), filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var sessions []models.WorkoutSession
	if err := cursor.All(context.TODO(), &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}
