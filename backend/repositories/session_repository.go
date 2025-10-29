package repositories

import (
	"context"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SessionRepositoryInterface interface {
	Create(ctx context.Context, session *models.Session) error
	FindActive(ctx context.Context) ([]models.Session, error)
	FindActiveByUser(ctx context.Context, userID primitive.ObjectID) ([]models.Session, error)
	MarkRevoked(ctx context.Context, id primitive.ObjectID, revokedAt time.Time, replacedByID *primitive.ObjectID) error
}

type SessionRepository struct {
	collection *mongo.Collection
}

func NewSessionRepository(db *mongo.Database) *SessionRepository {
	coll := db.Collection("sessions")
	createIndexes(coll)
	return &SessionRepository{collection: coll}
}

func createIndexes(coll *mongo.Collection) {
	_, _ = coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "userId", Value: 1}, {Key: "revokedAt", Value: 1}},
		},
		{
			Keys:    bson.D{{Key: "expiresAt", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	})
}
func (r *SessionRepository) FindActiveByUser(ctx context.Context, userID primitive.ObjectID) ([]models.Session, error) {
	return r.find(ctx, bson.M{
		"userId":    userID,
		"revokedAt": bson.M{"$exists": false},
	})
}

func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.collection.InsertOne(ctx, session)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		session.ID = oid
	}
	return nil
}

func (r *SessionRepository) FindActive(ctx context.Context) ([]models.Session, error) {
	return r.find(ctx, bson.M{
		"revokedAt": bson.M{"$exists": false},
	})
}

func (r *SessionRepository) find(ctx context.Context, filter bson.M) ([]models.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var sessions []models.Session
	if err := cursor.All(context.Background(), &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) MarkRevoked(ctx context.Context, id primitive.ObjectID, revokedAt time.Time, replacedByID *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"revokedAt": revokedAt,
		},
	}
	if replacedByID != nil {
		update["$set"].(bson.M)["replacedById"] = replacedByID
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
