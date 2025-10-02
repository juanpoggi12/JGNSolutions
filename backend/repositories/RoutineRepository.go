package repositories

import (
	"context"

	"github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoutineRepository struct {
	collection *mongo.Collection
}

func NewRoutineRepository(db *mongo.Database) *RoutineRepository {
	return &RoutineRepository{collection: db.Collection("routines")}
}

func (r *RoutineRepository) Create(ctx context.Context, routine *models.Routine) error {
	_, err := r.collection.InsertOne(ctx, routine)
	return err
}

func (r *RoutineRepository) FindByID(ctx context.Context, id string) (*models.Routine, error) {
	var routine models.Routine
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&routine)
	if err != nil {
		return nil, err
	}
	return &routine, nil
}

func (r *RoutineRepository) Update(ctx context.Context, routine *models.Routine) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": routine.ID}, bson.M{"$set": routine})
	return err
}

func (r *RoutineRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
