package repositories

import (
	"context"

	"github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WorkoutSessionRepository struct {
	collection *mongo.Collection
}

type WorkoutSessionRepositoryInterface interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.WorkoutSession, error)
	FindByUser(ctx context.Context, userID primitive.ObjectID) ([]models.WorkoutSession, error)
	Create(ctx context.Context, session *models.WorkoutSession) error
	Update(ctx context.Context, session *models.WorkoutSession) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

func NewWorkoutSessionRepository(db *mongo.Database) *WorkoutSessionRepository {
	return &WorkoutSessionRepository{collection: db.Collection("workoutSessions")}
}

func (r *WorkoutSessionRepository) Create(ctx context.Context, session *models.WorkoutSession) error {
	res, err := r.collection.InsertOne(ctx, session)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		session.ID = oid
	}

	return nil
}

func (r *WorkoutSessionRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.WorkoutSession, error) {
	var session models.WorkoutSession
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *WorkoutSessionRepository) Update(ctx context.Context, session *models.WorkoutSession) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": session.ID}, bson.M{"$set": session})
	return err
}

func (r *WorkoutSessionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *WorkoutSessionRepository) Search(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutSession, error) {
	cursor, err := r.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []models.WorkoutSession
	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

// Buscar sesiones por usuario
func (r *WorkoutSessionRepository) FindByUser(ctx context.Context, userID primitive.ObjectID) ([]models.WorkoutSession, error) {
	var sessions []models.WorkoutSession

	cursor, err := r.collection.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var s models.WorkoutSession
		if err := cursor.Decode(&s); err == nil {
			sessions = append(sessions, s)
		}
	}
	return sessions, cursor.Err()
}
