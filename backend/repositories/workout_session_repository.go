package repositories

import (
    "github.com/juanpoggi12/JGNSolutions/backend/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkoutRepository struct {
	collection *mongo.Collection
}

func NewWorkoutRepository(db *mongo.Database) *WorkoutRepository {
	return &WorkoutRepository{collection: db.Collection("workouts")}
}

func (r *WorkoutRepository) Create(ctx context.Context, workout *models.Workout) error {
	_, err := r.collection.InsertOne(ctx, workout)
	return err
}

func (r *WorkoutRepository) FindByID(ctx context.Context, id string) (*models.Workout, error) {
	var workout models.Workout
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&workout)
	if err != nil {
		return nil, err
	}
	return &workout, nil
}

func (r *WorkoutRepository) Update(ctx context.Context, workout *models.Workout) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": workout.ID}, bson.M{"$set": workout})
	return err
}

func (r *WorkoutRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
