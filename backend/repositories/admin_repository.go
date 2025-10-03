package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type AdminRepositoryInterface interface {
	ContarDocumentos(nombreColeccion string) (int64, error)
}

type AdminRepository struct {
	db *mongo.Database
}

func NewAdminRepository(db *mongo.Database) *AdminRepository {
	return &AdminRepository{db: db}
}

func (repository AdminRepository) ContarDocumentos(nombreColeccion string) (int64, error) {
	collection := repository.db.Collection(nombreColeccion)
	return collection.CountDocuments(context.TODO(), map[string]interface{}{})
}
