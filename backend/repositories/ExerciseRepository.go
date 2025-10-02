package repositories

import (
	"context"

	"github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExerciseRepository struct {
	collection *mongo.Collection
}

func NewExerciseRepository(db *mongo.Database) *ExerciseRepository {
	return &ExerciseRepository{collection: db.Collection("exercises")}
}

func (r *ExerciseRepository) Create(ctx context.Context, exercise *models.Exercise) error {
	_, err := r.collection.InsertOne(ctx, exercise)
	return err
}

func (r *ExerciseRepository) FindByID(ctx context.Context, id string) (*models.Exercise, error) {
	var exercise models.Exercise
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&exercise)
	if err != nil {
		return nil, err
	}
	return &exercise, nil
}

func (r *ExerciseRepository) Update(ctx context.Context, exercise *models.Exercise) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": exercise.ID}, bson.M{"$set": exercise})
	return err
}

func (r *ExerciseRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
