package repositories

import (
	"context"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogRepositoryInterface interface {
	InsertarLog(log models.Log) (*mongo.InsertOneResult, error)
	ListarLogs() ([]models.Log, error)
}

type LogRepository struct {
	db *mongo.Database
}

func NewLogRepository(db *mongo.Database) *LogRepository {
	return &LogRepository{db: db}
}

func (repository LogRepository) collection() *mongo.Collection {
	return repository.db.Collection("logs")
}

// InsertarLog → Guarda un nuevo registro en la colección de logs
func (repository LogRepository) InsertarLog(log models.Log) (*mongo.InsertOneResult, error) {
	collection := repository.collection()

	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	resultado, err := collection.InsertOne(context.TODO(), log)
	return resultado, err
}

// ListarLogs → Devuelve todos los logs ordenados por fecha descendente
func (repository LogRepository) ListarLogs() ([]models.Log, error) {
	collection := repository.collection()

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cur, err := collection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var logs []models.Log
	for cur.Next(context.Background()) {
		var log models.Log
		if err := cur.Decode(&log); err == nil {
			logs = append(logs, log)
		}
	}

	return logs, nil
}
