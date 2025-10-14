package repositories

import (
	"context"
	"time"

	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserProfileRepositoryInterface interface {
	ObtenerPerfilPorUserID(userID primitive.ObjectID) (models.UserProfile, error)
	InsertarPerfil(perfil models.UserProfile) (*mongo.InsertOneResult, error)
	ModificarPerfil(perfil models.UserProfile) (*mongo.UpdateResult, error)
	ListarPerfiles() ([]models.UserProfile, error)
	ObtenerPerfilPorID(id string) (models.UserProfile, error)
	EliminarPerfil(id primitive.ObjectID) (*mongo.DeleteResult, error)
}

type UserProfileRepository struct {
	db *mongo.Database
}

func NewUserProfileRepository(db *mongo.Database) *UserProfileRepository {
	return &UserProfileRepository{db: db}
}

func (repository UserProfileRepository) collection() *mongo.Collection {
	return repository.db.Collection("user_profiles")
}

// --- Métodos para usuarios ---

func (repository UserProfileRepository) ObtenerPerfilPorUserID(userID primitive.ObjectID) (models.UserProfile, error) {
	collection := repository.collection()
	filtro := bson.M{"userId": userID}

	var perfil models.UserProfile
	err := collection.FindOne(context.TODO(), filtro).Decode(&perfil)
	return perfil, err
}

func (repository UserProfileRepository) InsertarPerfil(perfil models.UserProfile) (*mongo.InsertOneResult, error) {
	collection := repository.collection()
	if perfil.UpdatedAt.IsZero() {
		perfil.UpdatedAt = time.Now()
	}
	resultado, err := collection.InsertOne(context.TODO(), perfil)
	return resultado, err
}

func (repository UserProfileRepository) ModificarPerfil(perfil models.UserProfile) (*mongo.UpdateResult, error) {
	collection := repository.collection()

	filtro := bson.M{"userId": perfil.UserID}
	actualizacion := bson.M{"$set": bson.M{
		"fullName":  perfil.FullName,
		"birthDate": perfil.BirthDate,
		"weightKg":  perfil.WeightKg,
		"heightCm":  perfil.HeightCm,
		"level":     perfil.Level,
		"goal":      perfil.Goal,
		"updatedAt": time.Now(),
	}}

	resultado, err := collection.UpdateOne(context.TODO(), filtro, actualizacion)
	return resultado, err
}

func (repository UserProfileRepository) ListarPerfiles() ([]models.UserProfile, error) {
	collection := repository.collection()
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var perfiles []models.UserProfile
	for cursor.Next(context.Background()) {
		var perfil models.UserProfile
		if err := cursor.Decode(&perfil); err == nil {
			perfiles = append(perfiles, perfil)
		}
	}
	return perfiles, nil
}

func (repository UserProfileRepository) ObtenerPerfilPorID(id string) (models.UserProfile, error) {
	collection := repository.collection()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.UserProfile{}, err
	}

	filtro := bson.M{"_id": objectID}
	var perfil models.UserProfile

	err = collection.FindOne(context.TODO(), filtro).Decode(&perfil)
	return perfil, err
}

func (repository UserProfileRepository) EliminarPerfil(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	collection := repository.collection()
	resultado, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return resultado, err
}
