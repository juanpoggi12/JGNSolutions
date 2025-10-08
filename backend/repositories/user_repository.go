package repositories

import (
	"context"

	"github.com/juanpoggi12/JGNSolutions/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepositoryInterface interface {
	InsertarUsuario(usuario models.User) (*mongo.InsertOneResult, error)
	ObtenerUsuarioPorID(id string) (models.User, error)
	ModificarUsuario(usuario models.User) (*mongo.UpdateResult, error)
	EliminarUsuario(id primitive.ObjectID) (*mongo.DeleteResult, error)
	ContarUsuarios() (int64, error)
	ListarUsuariosBasico() ([]models.User, error)
	FindByEmailOrUsername(identifier string) (*models.User, error)
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

// en el archivo (agregar implementación)
func (r UserRepository) ListarUsuariosBasico() ([]models.User, error) {
	coll := r.collection()

	// Proyección: solo _id, email y role
	opts := options.Find().SetProjection(bson.M{
		"_id": 1, "email": 1, "role": 1,
	})

	cur, err := coll.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var res []models.User
	for cur.Next(context.Background()) {
		var u models.User
		if err := cur.Decode(&u); err == nil {
			res = append(res, u)
		}
	}
	return res, nil
}

// Buscar por email o username (para login)
func (r *UserRepository) FindByEmailOrUsername(identifier string) (*models.User, error) {
	var user models.User
	filter := bson.M{
		"$or": []bson.M{
			{"email": identifier},
			{"username": identifier},
		},
	}
	err := r.db.Collection("users").FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
