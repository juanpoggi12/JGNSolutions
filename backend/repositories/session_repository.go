package repositories

import (
	"context"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionRepository struct {
	db *mongo.Database
}

func NewSessionRepository(db *mongo.Database) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(session models.Session) error {
	_, err := r.db.Collection("sessions").InsertOne(context.TODO(), session)
	return err
}

func (r *SessionRepository) FindValidByHash(hash string) (*models.Session, error) {
	var s models.Session
	err := r.db.Collection("sessions").FindOne(context.TODO(), bson.M{
		"refreshTokenHash": hash,
		"revokedAt":        bson.M{"$exists": false},
		"expiresAt":        bson.M{"$gt": time.Now()},
	}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SessionRepository) RevokeByID(id interface{}) error {
	now := time.Now()
	_, err := r.db.Collection("sessions").UpdateOne(context.TODO(),
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"revokedAt": now}},
	)
	return err
}
