package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// WorkoutEntryHandler maneja las operaciones sobre las entradas de entrenamiento
type WorkoutEntryHandler struct {
	entryService *services.WorkoutEntryService
}

// Constructor
func NewWorkoutEntryHandler(entryService *services.WorkoutEntryService) *WorkoutEntryHandler {
	return &WorkoutEntryHandler{entryService: entryService}
}

// POST /api/workout-entries → crear una nueva entrada
func (h *WorkoutEntryHandler) CreateEntry(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	var entry models.WorkoutEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.entryService.Create(ctx, actor, &entry); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

// GET /api/workout-entries/:id → obtener una entrada por ID
func (h *WorkoutEntryHandler) GetEntryByID(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	entry, err := h.entryService.GetByID(ctx, actor, id)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// PUT /api/workout-entries/:id → actualizar una entrada
func (h *WorkoutEntryHandler) UpdateEntry(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	var entry models.WorkoutEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	entry.ID = id

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.entryService.Update(ctx, actor, &entry); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entrada actualizada correctamente"})
}

// DELETE /api/workout-entries/:id → eliminar una entrada
func (h *WorkoutEntryHandler) DeleteEntry(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.entryService.Delete(ctx, actor, id); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entrada eliminada correctamente"})
}

// GET /api/workout-entries/search → buscar entradas
func (h *WorkoutEntryHandler) SearchEntries(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	// Filtros opcionales
	filter := bson.M{}
	if sessionID := c.Query("session_id"); sessionID != "" {
		if oid, err := primitive.ObjectIDFromHex(sessionID); err == nil {
			filter["workoutSessionId"] = oid
		}
	}
	if exerciseName := c.Query("exercise_name"); exerciseName != "" {
		filter["exerciseName"] = bson.M{"$regex": exerciseName, "$options": "i"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find()
	results, err := h.entryService.Search(ctx, actor, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
