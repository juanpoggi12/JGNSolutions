package repositories

import (
	"context"

	"github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryInterface interface {
	InsertarUsuario(usuario models.User) (*mongo.InsertOneResult, error)
	ObtenerUsuarioPorID(id string) (models.User, error)
	ModificarUsuario(usuario models.User) (*mongo.UpdateResult, error)
	EliminarUsuario(id primitive.ObjectID) (*mongo.DeleteResult, error)
	ContarUsuarios() (int64, error)
}

type UserRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (repository UserRepository) collection() *mongo.Collection {
	return repository.db.Collection("users")
}

func (repository UserRepository) InsertarUsuario(usuario models.User) (*mongo.InsertOneResult, error) {
	collection := repository.collection()
	resultado, err := collection.InsertOne(context.TODO(), usuario)
	return resultado, err
}

func (repository UserRepository) ObtenerUsuarioPorID(id string) (models.User, error) {
	collection := repository.collection()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, err
	}

	filtro := bson.M{"_id": objectID}
	var usuario models.User

	err = collection.FindOne(context.TODO(), filtro).Decode(&usuario)
	return usuario, err
}

func (repository UserRepository) ModificarUsuario(usuario models.User) (*mongo.UpdateResult, error) {
	collection := repository.collection()

	filtro := bson.M{"_id": usuario.ID}
	actualizacion := bson.M{"$set": bson.M{
		"username":     usuario.Username,
		"email":        usuario.Email,
		"passwordHash": usuario.PasswordHash,
		"role":         usuario.Role,
		"createdAt":    usuario.CreatedAt,
		"updatedAt":    usuario.UpdatedAt,
		"isActive":     usuario.IsActive,
	}}

	resultado, err := collection.UpdateOne(context.TODO(), filtro, actualizacion)
	return resultado, err
}

func (repository UserRepository) EliminarUsuario(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	collection := repository.collection()
	resultado, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return resultado, err
}

func (repository UserRepository) ContarUsuarios() (int64, error) {
	collection := repository.collection()
	return collection.CountDocuments(context.TODO(), bson.M{})
}
