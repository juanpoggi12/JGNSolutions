package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureExerciseIndexes(ctx context.Context, coll *mongo.Collection, uniquePerCreator bool) error {
	models := []mongo.IndexModel{
		{Keys: bson.D{{Key: "categoria", Value: 1}}},
		{Keys: bson.D{{Key: "grupoMuscular", Value: 1}}},
		{Keys: bson.D{{Key: "createdByUserId", Value: 1}}},
		{Keys: bson.D{{Key: "isDeleted", Value: 1}}}, // Ãºtil p/filtrar soft delete
	}

	if uniquePerCreator {
		models = append(models, mongo.IndexModel{
			Keys:    bson.D{{Key: "createdByUserId", Value: 1}, {Key: "nombre", Value: 1}},
			Options: options.Index().SetUnique(true),
		})
	} else {
		models = append(models, mongo.IndexModel{
			Keys:    bson.D{{Key: "nombre", Value: 1}},
			Options: options.Index().SetUnique(true),
		})
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	return err
}
func EnsureRoutineIndexes(ctx context.Context, coll *mongo.Collection) error {
	models := []mongo.IndexModel{
		{Keys: bson.D{{Key: "userId", Value: 1}}},
		{Keys: bson.D{{Key: "nombre", Value: 1}}},
		{Keys: bson.D{{Key: "isDeleted", Value: 1}}},
	}
	_, err := coll.Indexes().CreateMany(ctx, models)
	return err
}
func EnsureRoutineExerciseIndexes(ctx context.Context, coll *mongo.Collection) error {
	models := []mongo.IndexModel{
		{Keys: bson.D{{Key: "routineId", Value: 1}}},
		{Keys: bson.D{{Key: "exerciseId", Value: 1}}},
		{Keys: bson.D{{Key: "routineId", Value: 1}, {Key: "orden", Value: 1}}},
	}
	_, err := coll.Indexes().CreateMany(ctx, models)
	return err
}
func EnsureWorkoutSessionIndexes(ctx context.Context, coll *mongo.Collection) error {
	models := []mongo.IndexModel{
		{Keys: bson.D{{Key: "userId", Value: 1}}},
		{Keys: bson.D{{Key: "fechaHoraInicio", Value: -1}}},
	}
	_, err := coll.Indexes().CreateMany(ctx, models)
	return err
}
func EnsureWorkoutEntryIndexes(ctx context.Context, coll *mongo.Collection) error {
	models := []mongo.IndexModel{
		{Keys: bson.D{{Key: "workoutSessionId", Value: 1}}},
		{Keys: bson.D{{Key: "exerciseId", Value: 1}}},
	}
	_, err := coll.Indexes().CreateMany(ctx, models)
	return err
}

func EnsureSessionIndexes(ctx context.Context, coll *mongo.Collection) error {
	models := []mongo.IndexModel{
		{Keys: bson.D{{Key: "userId", Value: 1}}},
		{Keys: bson.D{{Key: "expiresAt", Value: 1}}},
	}
	_, err := coll.Indexes().CreateMany(ctx, models)
	return err
}
