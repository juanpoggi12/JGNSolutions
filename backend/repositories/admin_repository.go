package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Struct para el resultado de la agregación de ejercicios más usados
type ExerciseUsageAgg struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"name"`
	UsageCount int                `bson:"usageCount"`
}

// Struct para el resultado de la agregación de rutinas más usadas
type RoutineUsageAgg struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"name"`
	UsageCount int                `bson:"usageCount"`
}

type AdminRepositoryInterface interface {
	CountDocuments(collectionName string) (int64, error)
	TopExercisesByEntries(limit int) ([]ExerciseUsageAgg, error)
	TopRoutinesBySessions(limit int) ([]RoutineUsageAgg, error)
}

type AdminRepository struct {
	db *mongo.Database
}

func NewAdminRepository(db *mongo.Database) *AdminRepository {
	return &AdminRepository{db: db}
}

func (repository AdminRepository) CountDocuments(collectionName string) (int64, error) {
	collection := repository.db.Collection(collectionName)
	return collection.CountDocuments(context.TODO(), map[string]interface{}{})
}

func (r AdminRepository) TopExercisesByEntries(limit int) ([]ExerciseUsageAgg, error) {
	if limit <= 0 {
		limit = 10
	}
	coll := r.db.Collection("workoutEntries")

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$exerciseId"},
			{Key: "usageCount", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "usageCount", Value: -1}}}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "exercises"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "ex"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$ex"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "usageCount", Value: 1},
			{Key: "name", Value: bson.D{{Key: "$ifNull", Value: bson.A{"$ex.name", "N/A"}}}},
		}}},
	}

	cur, err := coll.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var out []ExerciseUsageAgg
	for cur.Next(context.Background()) {
		var row ExerciseUsageAgg
		if err := cur.Decode(&row); err == nil {
			out = append(out, row)
		}
	}
	return out, nil
}

func (r AdminRepository) TopRoutinesBySessions(limit int) ([]RoutineUsageAgg, error) {
	if limit <= 0 {
		limit = 10
	}
	coll := r.db.Collection("workoutSessions")

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "routineId", Value: bson.D{{Key: "$type", Value: "objectId"}}}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$routineId"},
			{Key: "usageCount", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "usageCount", Value: -1}}}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "routines"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "r"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$r"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "usageCount", Value: 1},
			{Key: "name", Value: bson.D{{Key: "$ifNull", Value: bson.A{"$r.name", "N/A"}}}},
		}}},
	}

	cur, err := coll.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var out []RoutineUsageAgg
	for cur.Next(context.Background()) {
		var row RoutineUsageAgg
		if err := cur.Decode(&row); err == nil {
			out = append(out, row)
		}
	}
	return out, nil
}
