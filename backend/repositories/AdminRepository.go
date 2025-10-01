package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type AdminRepository struct {
	db *mongo.Database
}

func NewAdminRepository(db *mongo.Database) *AdminRepository {
	return &AdminRepository{db: db}
}

// Ejemplo de método genérico
func (r *AdminRepository) CountCollection(ctx context.Context, name string) (int, error) {
	collection := r.db.Collection(name)
	count, err := collection.CountDocuments(ctx, map[string]interface{}{})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
