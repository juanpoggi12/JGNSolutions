package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func Connect(ctx context.Context, uri, dbName string) (*Mongo, error) {
	cl, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := cl.Ping(ctxPing, nil); err != nil {
		return nil, err
	}
	return &Mongo{Client: cl, DB: cl.Database(dbName)}, nil
}
