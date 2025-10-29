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
	GetProfileByUserID(userID primitive.ObjectID) (models.UserProfile, error)
	UpdateProfile(profile models.UserProfile) (*mongo.UpdateResult, error)
	ListProfiles() ([]models.UserProfile, error)
	GetProfileByID(id string) (models.UserProfile, error)
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

// Create inserts a new user profile (used by AuthService to create default profiles)
func (repository *UserProfileRepository) Create(ctx context.Context, profile *models.UserProfile) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if profile.UpdatedAt.IsZero() {
		profile.UpdatedAt = time.Now()
	}

	res, err := repository.collection().InsertOne(ctx, profile)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		profile.ID = oid
	}
	return nil
}

// --- Métodos para usuarios ---

func (repository UserProfileRepository) GetProfileByUserID(userID primitive.ObjectID) (models.UserProfile, error) {
	collection := repository.collection()
	filtro := bson.M{"userId": userID}

	var profile models.UserProfile
	err := collection.FindOne(context.TODO(), filtro).Decode(&profile)
	return profile, err
}

func (repository UserProfileRepository) UpdateProfile(profile models.UserProfile) (*mongo.UpdateResult, error) {
	collection := repository.collection()

	filtro := bson.M{"userId": profile.UserID}
	actualizacion := bson.M{"$set": bson.M{
		"fullName":  profile.FullName,
		"birthDate": profile.BirthDate,
		"weightKg":  profile.WeightKg,
		"heightCm":  profile.HeightCm,
		"level":     profile.Level,
		"goal":      profile.Goal,
		"updatedAt": time.Now(),
	}}

	resultado, err := collection.UpdateOne(context.TODO(), filtro, actualizacion)
	return resultado, err
}

func (repository UserProfileRepository) ListProfiles() ([]models.UserProfile, error) {
	collection := repository.collection()
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var profiles []models.UserProfile
	for cursor.Next(context.Background()) {
		var profile models.UserProfile
		if err := cursor.Decode(&profile); err == nil {
			profiles = append(profiles, profile)
		}
	}
	return profiles, nil
}

func (repository UserProfileRepository) GetProfileByID(id string) (models.UserProfile, error) {
	collection := repository.collection()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.UserProfile{}, err
	}

	filtro := bson.M{"_id": objectID}
	var profile models.UserProfile

	err = collection.FindOne(context.TODO(), filtro).Decode(&profile)
	return profile, err
}
