package repositories

import (
	"context"

	"github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WorkoutEntryRepository struct {
	collection *mongo.Collection
}

func NewWorkoutEntryRepository(db *mongo.Database) *WorkoutEntryRepository {
	return &WorkoutEntryRepository{collection: db.Collection("workoutEntries")}
}

func (r *WorkoutEntryRepository) Create(ctx context.Context, entry *models.WorkoutEntry) error {
	res, err := r.collection.InsertOne(ctx, entry)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		entry.ID = oid
	}

	return nil
}

func (r *WorkoutEntryRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.WorkoutEntry, error) {
	var entry models.WorkoutEntry
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *WorkoutEntryRepository) Update(ctx context.Context, entry *models.WorkoutEntry) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": entry.ID}, bson.M{"$set": entry})
	return err
}

func (r *WorkoutEntryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *WorkoutEntryRepository) Search(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutEntry, error) {
	cursor, err := r.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var entries []models.WorkoutEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}
