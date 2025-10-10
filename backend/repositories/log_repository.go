package repositories

import (
	"context"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogRepository struct {
	db *mongo.Database
}

func NewLogRepository(db *mongo.Database) *LogRepository {
	return &LogRepository{db: db}
}

func (r LogRepository) collection() *mongo.Collection {
	return r.db.Collection("logs")
}

// Insertar un nuevo log
func (r LogRepository) InsertarLog(ctx context.Context, log *models.Log) (*mongo.InsertOneResult, error) {
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}
	return r.collection().InsertOne(ctx, log)
}

// Listar todos los logs ordenados por fecha (más recientes primero)
func (r LogRepository) ListarLogs(ctx context.Context) ([]models.Log, error) {
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cur, err := r.collection().Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var logs []models.Log
	for cur.Next(ctx) {
		var l models.Log
		if err := cur.Decode(&l); err == nil {
			logs = append(logs, l)
		}
	}
	return logs, cur.Err()
}
