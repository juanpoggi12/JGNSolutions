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

type UserStatsRepositoryInterface interface {
	FindSessions(ctx context.Context, userID primitive.ObjectID, from, to *time.Time) ([]models.WorkoutSession, error)
	FindEntriesBySessions(ctx context.Context, sessionIDs []primitive.ObjectID) ([]models.WorkoutEntry, error)
	FindExercisesByIDs(ctx context.Context, ids []primitive.ObjectID) (map[primitive.ObjectID]models.Exercise, error)
	FindRoutinesByIDs(ctx context.Context, ids []primitive.ObjectID) (map[primitive.ObjectID]models.Routine, error)
}

type UserStatsRepository struct {
	db *mongo.Database
}

func NewUserStatsRepository(db *mongo.Database) *UserStatsRepository {
	return &UserStatsRepository{db: db}
}

func (r *UserStatsRepository) sessionsCollection() *mongo.Collection {
	return r.db.Collection("workoutSessions")
}

func (r *UserStatsRepository) entriesCollection() *mongo.Collection {
	return r.db.Collection("workoutEntries")
}

func (r *UserStatsRepository) exercisesCollection() *mongo.Collection {
	return r.db.Collection("exercises")
}

func (r *UserStatsRepository) routinesCollection() *mongo.Collection {
	return r.db.Collection("routines")
}

func (r *UserStatsRepository) FindSessions(ctx context.Context, userID primitive.ObjectID, from, to *time.Time) ([]models.WorkoutSession, error) {
	if userID.IsZero() {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"userId": userID}
	if from != nil || to != nil {
		timeFilter := bson.M{}
		if from != nil {
			timeFilter["$gte"] = *from
		}
		if to != nil {
			timeFilter["$lt"] = *to
		}
		filter["startTime"] = timeFilter
	}

	opts := options.Find().SetSort(bson.M{"startTime": 1})

	cursor, err := r.sessionsCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var sessions []models.WorkoutSession
	if err := cursor.All(context.Background(), &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *UserStatsRepository) FindEntriesBySessions(ctx context.Context, sessionIDs []primitive.ObjectID) ([]models.WorkoutEntry, error) {
	if len(sessionIDs) == 0 {
		return []models.WorkoutEntry{}, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"workoutSessionId": bson.M{"$in": sessionIDs}}
	cursor, err := r.entriesCollection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var entries []models.WorkoutEntry
	if err := cursor.All(context.Background(), &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *UserStatsRepository) FindExercisesByIDs(ctx context.Context, ids []primitive.ObjectID) (map[primitive.ObjectID]models.Exercise, error) {
	result := make(map[primitive.ObjectID]models.Exercise)
	if len(ids) == 0 {
		return result, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := r.exercisesCollection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var ex models.Exercise
		if err := cursor.Decode(&ex); err == nil {
			result[ex.ID] = ex
		}
	}

	return result, nil
}

func (r *UserStatsRepository) FindRoutinesByIDs(ctx context.Context, ids []primitive.ObjectID) (map[primitive.ObjectID]models.Routine, error) {
	result := make(map[primitive.ObjectID]models.Routine)
	if len(ids) == 0 {
		return result, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := r.routinesCollection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var routine models.Routine
		if err := cursor.Decode(&routine); err == nil {
			result[routine.ID] = routine
		}
	}

	return result, nil
}
