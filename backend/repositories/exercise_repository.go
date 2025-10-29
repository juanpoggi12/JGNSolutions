package repositories

import (
	"context"
	"strings"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ExerciseRepositoryInterface interface {
	GetExerciseByID(id string) (models.Exercise, error)
	InsertExercise(exercise models.Exercise) (*mongo.InsertOneResult, error)
	UpdateExercise(exercise models.Exercise) (*mongo.UpdateResult, error)
	DeleteExercise(id primitive.ObjectID) (*mongo.UpdateResult, error)
	ListExercisesCatalog(ctx context.Context, params ExerciseCatalogParams) ([]models.Exercise, int64, error)
}

type ExerciseRepository struct {
	db *mongo.Database
}

type ExerciseCatalogParams struct {
	Query          string
	Category       string
	MuscleGroup    string
	Difficulty     string
	IncludeDeleted bool
	Page           int
	Limit          int
}

func NewExerciseRepository(db *mongo.Database) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

func (repository ExerciseRepository) collection() *mongo.Collection {
	return repository.db.Collection("exercises")
}

func (repository ExerciseRepository) GetExerciseByID(id string) (models.Exercise, error) {
	collection := repository.collection()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Exercise{}, err
	}

	filtro := bson.M{"_id": objectID}
	var exercise models.Exercise

	err = collection.FindOne(context.TODO(), filtro).Decode(&exercise)
	return exercise, err
}

func (repository ExerciseRepository) InsertExercise(exercise models.Exercise) (*mongo.InsertOneResult, error) {
	collection := repository.collection()
	resultado, err := collection.InsertOne(context.TODO(), exercise)
	return resultado, err
}

func (repository ExerciseRepository) UpdateExercise(exercise models.Exercise) (*mongo.UpdateResult, error) {
	collection := repository.collection()

	filtro := bson.M{"_id": exercise.ID}
	actualizacion := bson.M{"$set": bson.M{
		"name":         exercise.Name,
		"description":  exercise.Description,
		"category":     exercise.Category,
		"muscleGroup":  exercise.MuscleGroup,
		"difficulty":   exercise.Difficulty,
		"mediaUrl":     exercise.MediaURL,
		"instructions": exercise.Instructions,
		"createdBy":    exercise.CreatedBy,
		"updatedAt":    exercise.UpdatedAt,
	}}

	resultado, err := collection.UpdateOne(context.TODO(), filtro, actualizacion)
	return resultado, err
}

func (repository ExerciseRepository) DeleteExercise(id primitive.ObjectID) (*mongo.UpdateResult, error) {
	collection := repository.collection()

	filtro := bson.M{"_id": id}
	actualizacion := bson.M{"$set": bson.M{
		"isDeleted": true,
		"updatedAt": time.Now(),
	}}

	resultado, err := collection.UpdateOne(context.TODO(), filtro, actualizacion)
	return resultado, err
}
func (repository ExerciseRepository) ListExercisesCatalog(ctx context.Context, params ExerciseCatalogParams) ([]models.Exercise, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if params.Query != "" {
		filter["name"] = bson.M{"$regex": params.Query, "$options": "i"}
	}
	if params.Category != "" {
		filter["category"] = strings.ToUpper(params.Category)
	}
	if params.MuscleGroup != "" {
		filter["muscleGroup"] = strings.ToUpper(params.MuscleGroup)
	}
	if params.Difficulty != "" {
		filter["difficulty"] = strings.ToUpper(params.Difficulty)
	}
	if !params.IncludeDeleted {
		filter["isDeleted"] = false
	}

	collection := repository.collection()

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	page := params.Page
	if page <= 0 {
		page = 1
	}
	limit := params.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	skip := int64((page - 1) * limit)
	opts := options.Find().SetSkip(skip).SetLimit(int64(limit)).SetSort(bson.M{"name": 1})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	var ejercicios []models.Exercise
	if err := cursor.All(context.Background(), &ejercicios); err != nil {
		return nil, 0, err
	}

	return ejercicios, total, nil
}
