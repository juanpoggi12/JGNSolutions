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

type UserRepositoryInterface interface {
	InsertUser(user models.User) (*mongo.InsertOneResult, error)
	GetUserByID(id string) (models.User, error)
	UpdateUser(user models.User) (*mongo.UpdateResult, error)
	DeleteUser(id primitive.ObjectID) (*mongo.DeleteResult, error)
	CountUsers() (int64, error)
	ListUsersBasic() ([]models.User, error)
	FindByEmailOrUsername(identifier string) (*models.User, error)
	UpdatePassword(id primitive.ObjectID, hashed string) (*mongo.UpdateResult, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
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

func (repository UserRepository) InsertUser(user models.User) (*mongo.InsertOneResult, error) {
	collection := repository.collection()
	resultado, err := collection.InsertOne(context.TODO(), user)
	return resultado, err
}

func (repository UserRepository) GetUserByID(id string) (models.User, error) {
	collection := repository.collection()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, err
	}

	filtro := bson.M{"_id": objectID}
	var user models.User

	err = collection.FindOne(context.TODO(), filtro).Decode(&user)
	return user, err
}

func (repository UserRepository) UpdateUser(user models.User) (*mongo.UpdateResult, error) {
	collection := repository.collection()

	filtro := bson.M{"_id": user.ID}
	actualizacion := bson.M{"$set": bson.M{
		"username":     user.Username,
		"email":        user.Email,
		"passwordHash": user.PasswordHash,
		"role":         user.Role,
		"createdAt":    user.CreatedAt,
		"updatedAt":    user.UpdatedAt,
		"isActive":     user.IsActive,
	}}

	resultado, err := collection.UpdateOne(context.TODO(), filtro, actualizacion)
	return resultado, err
}

func (repository UserRepository) DeleteUser(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	collection := repository.collection()
	resultado, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return resultado, err
}

func (repository UserRepository) CountUsers() (int64, error) {
	collection := repository.collection()
	return collection.CountDocuments(context.TODO(), bson.M{})
}

// en el archivo (agregar implementación)
func (r UserRepository) ListUsersBasic() ([]models.User, error) {
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

func (repository UserRepository) UpdatePassword(id primitive.ObjectID, hashed string) (*mongo.UpdateResult, error) {
	// ... implementación existente ...
	collection := repository.collection()
	filtro := bson.M{"_id": id}
	actualizacion := bson.M{
		"$set": bson.M{
			"passwordHash": hashed,
			"updatedAt":    time.Now(),
		},
	}
	return collection.UpdateOne(context.TODO(), filtro, actualizacion)
}
func (repository *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	count, err := repository.collection().CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repository *UserRepository) Create(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := repository.collection().InsertOne(ctx, user)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	}
	return nil
}

func (repository *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user models.User
	if err := repository.collection().FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repository *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user models.User
	if err := repository.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
