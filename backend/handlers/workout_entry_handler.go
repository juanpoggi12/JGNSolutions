package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/services"
	"github.com/juanpoggi12/JGNSolutions/backend/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// WorkoutEntryHandler maneja las operaciones sobre las entradas de entrenamiento
type WorkoutEntryHandler struct {
	entryService services.WorkoutEntryServiceInterface
}

// Constructor
func NewWorkoutEntryHandler(entryService services.WorkoutEntryServiceInterface) *WorkoutEntryHandler {
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

	var req dto.WorkoutEntryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	created, err := h.entryService.Create(actor, req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	responseDTO := utils.ConvertWorkoutEntryModelToResponse(created)
	c.JSON(http.StatusCreated, responseDTO)

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

	entry, err := h.entryService.GetByID(actor, id)
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

	var req dto.WorkoutEntryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	updated, err := h.entryService.Update(actor, id, req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
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

	if err := h.entryService.Delete(actor, id); err != nil {
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

	filter := bson.M{}
	if sessionID := c.Query("session_id"); sessionID != "" {
		if oid, err := primitive.ObjectIDFromHex(sessionID); err == nil {
			filter["workoutSessionId"] = oid
		}
	}

	if exerciseID := c.Query("exercise_id"); exerciseID != "" {
		if oid, err := primitive.ObjectIDFromHex(exerciseID); err == nil {
			filter["exerciseId"] = oid
		}
	}

	opts := options.Find()
	results, err := h.entryService.Search(actor, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
