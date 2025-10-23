// juanpoggi12/jgnsolutions/JGNSolutions-7c48b53190321ccfabc6877d44ae535f756457c5/backend/repositories/workout_session_repository.go
package repositories

import (
	"context"
	"errors" // Import errors package
	"fmt"    // Import fmt package
	"log"    // Import log package
	"time"   // Import time package

	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Interfaz del repositorio (contrato que usa el service)
type WorkoutSessionRepositoryInterface interface {
	Create(session *models.WorkoutSession) error
	FindByID(id primitive.ObjectID) (*models.WorkoutSession, error)
	FindByUser(userID primitive.ObjectID) ([]models.WorkoutSession, error)
	Update(session *models.WorkoutSession) error
	Delete(id primitive.ObjectID) error
	Search(filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutSession, error)
}

// Implementación concreta (usa MongoDB)
type WorkoutSessionRepository struct {
	collection *mongo.Collection
}

// Constructor
func NewWorkoutSessionRepository(db *mongo.Database) *WorkoutSessionRepository {
	return &WorkoutSessionRepository{collection: db.Collection("workoutSessions")}
}

// Crear una sesión
func (r *WorkoutSessionRepository) Create(session *models.WorkoutSession) error {
	// Ensure CreatedAt and UpdatedAt are set if zero
	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now()
	}
	if session.UpdatedAt.IsZero() {
		session.UpdatedAt = time.Now()
	}
	// Explicitly set a new ObjectID if the session doesn't have one
	// This helps prevent accidental updates if an existing object is passed
	if session.ID.IsZero() {
		session.ID = primitive.NewObjectID()
		log.Printf("[Repo.Create] Generated new ObjectID for session: %s", session.ID.Hex())
	} else {
		// If an ID already exists, log a warning - this function should ideally receive a new object
		log.Printf("[Repo.Create] WARNING: Session object passed to Create already has an ID: %s. Attempting InsertOne anyway.", session.ID.Hex())
		// Consider returning an error here if Create should *only* insert new documents
		// return errors.New("cannot create session: object already has an ID")
	}

	log.Printf("[Repo.Create] Attempting InsertOne for session with tentative ID: %s", session.ID.Hex())
	res, err := r.collection.InsertOne(context.TODO(), session) // Pass the session object directly
	if err != nil {
		log.Printf("[Repo.Create] Error during InsertOne: %v", err)
		// Check for duplicate key error specifically (if you added a unique index)
		if mongo.IsDuplicateKeyError(err) {
			log.Printf("[Repo.Create] Duplicate key error. Session ID %s might already exist.", session.ID.Hex())
			return fmt.Errorf("error al insertar sesión: ID duplicado (%s)", session.ID.Hex())
		}
		return fmt.Errorf("error al insertar la sesión en la base de datos: %w", err)
	}

	// Double-check the returned InsertedID matches the ID we tried to insert
	returnedOID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("[Repo.Create] InsertOne succeeded but InsertedID is not an ObjectID (%T): %v. Session object ID remains %s.", res.InsertedID, res.InsertedID, session.ID.Hex())
		// Stick with the ID assigned *before* InsertOne: session.ID
	} else if returnedOID != session.ID {
		log.Printf("[Repo.Create] WARNING: InsertOne returned ObjectID %s which DIFFERS from the pre-generated ID %s. Using the pre-generated ID.", returnedOID.Hex(), session.ID.Hex())
		// Stick with the ID assigned *before* InsertOne: session.ID
	} else {
		log.Printf("[Repo.Create] InsertOne successful. Returned InsertedID %s matches pre-generated ID %s.", returnedOID.Hex(), session.ID.Hex())
	}

	// Log the final ID value *in the session object* before returning
	log.Printf("[Repo.Create] Session object ID before returning: %s", session.ID.Hex())
	return nil // Success
}

// Buscar por ID
func (r *WorkoutSessionRepository) FindByID(id primitive.ObjectID) (*models.WorkoutSession, error) {
	var session models.WorkoutSession
	if err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&session); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("sesión no encontrada con ese ID") // Más específico
		}
		log.Printf("[Repo.FindByID] Error finding session %s: %v", id.Hex(), err)
		return nil, err // Otro error
	}
	return &session, nil
}

// Buscar todas las sesiones de un usuario
func (r *WorkoutSessionRepository) FindByUser(userID primitive.ObjectID) ([]models.WorkoutSession, error) {
	var sessions []models.WorkoutSession

	// Añadir ordenación por defecto, por ejemplo, más recientes primero
	opts := options.Find().SetSort(bson.D{{Key: "startTime", Value: -1}})

	cursor, err := r.collection.Find(context.TODO(), bson.M{"userId": userID}, opts)
	if err != nil {
		log.Printf("[Repo.FindByUser] Error finding sessions for user %s: %v", userID.Hex(), err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &sessions); err != nil {
		log.Printf("[Repo.FindByUser] Error decoding sessions for user %s: %v", userID.Hex(), err)
		return nil, err
	}

	return sessions, nil
}

// Actualizar una sesión
func (r *WorkoutSessionRepository) Update(session *models.WorkoutSession) error {
	// Asegura que UpdatedAt se actualice
	session.UpdatedAt = time.Now()
	// Usa UpdateOne con $set para actualizar solo los campos necesarios
	// Asume que 'session' tiene los campos actualizados
	update := bson.M{"$set": session}
	log.Printf("[Repo.Update] Attempting UpdateOne for session ID: %s", session.ID.Hex())
	res, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": session.ID}, update)
	if err != nil {
		log.Printf("[Repo.Update] Error updating session %s: %v", session.ID.Hex(), err)
		return err
	}
	if res.MatchedCount == 0 {
		log.Printf("[Repo.Update] No session found with ID %s to update.", session.ID.Hex())
		return mongo.ErrNoDocuments // O un error personalizado
	}
	log.Printf("[Repo.Update] Successfully updated session %s (Matched: %d, Modified: %d)", session.ID.Hex(), res.MatchedCount, res.ModifiedCount)
	return nil
}

// Eliminar una sesión
func (r *WorkoutSessionRepository) Delete(id primitive.ObjectID) error {
	log.Printf("[Repo.Delete] Attempting DeleteOne for session ID: %s", id.Hex())
	res, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		log.Printf("[Repo.Delete] Error deleting session %s: %v", id.Hex(), err)
		return err
	}
	if res.DeletedCount == 0 {
		log.Printf("[Repo.Delete] No session found with ID %s to delete.", id.Hex())
		return mongo.ErrNoDocuments // O error personalizado
	}
	log.Printf("[Repo.Delete] Successfully deleted session %s", id.Hex())
	return nil
}

// Buscar con filtros personalizados
func (r *WorkoutSessionRepository) Search(filter bson.M, opts ...*options.FindOptions) ([]models.WorkoutSession, error) {
	// Añade ordenación por defecto si no se proporciona
	finalOpts := options.MergeFindOptions(opts...)
	if finalOpts.Sort == nil {
		finalOpts.SetSort(bson.D{{Key: "startTime", Value: -1}}) // Más recientes primero
	}

	log.Printf("[Repo.Search] Executing Find with filter: %v", filter)
	cursor, err := r.collection.Find(context.TODO(), filter, finalOpts)
	if err != nil {
		log.Printf("[Repo.Search] Error executing find: %v", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var sessions []models.WorkoutSession
	if err := cursor.All(context.TODO(), &sessions); err != nil {
		log.Printf("[Repo.Search] Error decoding results: %v", err)
		return nil, err
	}
	log.Printf("[Repo.Search] Found %d sessions.", len(sessions))
	return sessions, nil
}
