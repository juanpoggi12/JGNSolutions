package handlers

import (
	"log" // Import log package
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/services"
	"github.com/juanpoggi12/JGNSolutions/backend/utils"

	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// WorkoutSessionHandler maneja las operaciones sobre sesiones de entrenamiento
type WorkoutSessionHandler struct {
	sessionService services.WorkoutSessionServiceInterface
}

// Constructor
func NewWorkoutSessionHandler(sessionService services.WorkoutSessionServiceInterface) *WorkoutSessionHandler {
	return &WorkoutSessionHandler{sessionService: sessionService}
}

// POST /api/workout-sessions → crear una nueva sesión (user o admin)
func (h *WorkoutSessionHandler) CreateSession(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}
	log.Printf("[Handler.CreateSession] Received request from actor %s (role %s)", actor.UserID.Hex(), actor.Role)

	var req dto.WorkoutSessionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Handler.CreateSession] Invalid JSON data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}
	log.Printf("[Handler.CreateSession] Request data bound successfully.")

	createdModel, err := h.sessionService.Create(actor, req)
	if err != nil {
		log.Printf("[Handler.CreateSession] Error received from sessionService.Create: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la sesión de entrenamiento"})
		return
	}

	if createdModel.ID.IsZero() {

		log.Printf("[Handler.CreateSession] CRITICAL ERROR: Service returned success, but WorkoutSession model ID is Zero!")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al generar ID de sesión"})
		return
	}
	log.Printf("[Handler.CreateSession] Service call successful. Model ID received: %s", createdModel.ID.Hex())

	responseDTO := utils.ConvertWorkoutSessionModelToResponse(createdModel)
	log.Printf("[Handler.CreateSession] Converted model to response DTO. DTO ID: %s", responseDTO.ID)

	log.Printf("[Handler.CreateSession] Sending StatusCreated with DTO.")
	c.JSON(http.StatusCreated, responseDTO)
}

// GET /api/workout-sessions/:id → obtener una sesión
func (h *WorkoutSessionHandler) GetSessionByID(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		log.Printf("[Handler.GetSessionByID] Invalid ID format: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de sesión inválido"})
		return
	}
	log.Printf("[Handler.GetSessionByID] Request for session %s by actor %s", id.Hex(), actor.UserID.Hex())

	sessionModel, err := h.sessionService.GetByID(actor, id)
	if err != nil {
		log.Printf("[Handler.GetSessionByID] Error from service GetByID for session %s: %v", id.Hex(), err)
		// Handle specific errors from service (e.g., not found, permission denied)
		if err.Error() == "sesión no encontrada" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err.Error() == "no tienes permiso para acceder a esta sesión" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener la sesión"})
		}
		return
	}

	// Convert model to DTO before sending
	responseDTO := utils.ConvertWorkoutSessionModelToResponse(*sessionModel)
	log.Printf("[Handler.GetSessionByID] Session %s found and authorized. Returning DTO.", id.Hex())
	c.JSON(http.StatusOK, responseDTO)
}

// PUT /api/workout-sessions/:id → actualizar sesión
func (h *WorkoutSessionHandler) UpdateSession(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		log.Printf("[Handler.UpdateSession] Invalid ID format: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de sesión inválido"})
		return
	}
	log.Printf("[Handler.UpdateSession] Request to update session %s by actor %s", id.Hex(), actor.UserID.Hex())

	var req dto.WorkoutSessionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Handler.UpdateSession] Invalid JSON data for session %s: %v", id.Hex(), err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de actualización inválidos: " + err.Error()})
		return
	}
	log.Printf("[Handler.UpdateSession] Update request data bound successfully for session %s.", id.Hex())

	updatedModel, err := h.sessionService.Update(actor, id, req)
	if err != nil {
		log.Printf("[Handler.UpdateSession] Error from service Update for session %s: %v", id.Hex(), err)
		if err.Error() == "sesión no encontrada para actualizar" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sesión no encontrada"})
		} else if err.Error() == "no tienes permiso para modificar esta sesión" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else if err.Error() == "error en los datos de actualización" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Propagate bad request from service
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la sesión"})
		}
		return
	}

	// Convert updated model to DTO
	responseDTO := utils.ConvertWorkoutSessionModelToResponse(updatedModel)
	log.Printf("[Handler.UpdateSession] Session %s updated successfully. Returning DTO.", id.Hex())
	c.JSON(http.StatusOK, responseDTO)
}

// DELETE /api/workout-sessions/:id → eliminar sesión
func (h *WorkoutSessionHandler) DeleteSession(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		log.Printf("[Handler.DeleteSession] Invalid ID format: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de sesión inválido"})
		return
	}
	log.Printf("[Handler.DeleteSession] Request to delete session %s by actor %s", id.Hex(), actor.UserID.Hex())

	if err := h.sessionService.Delete(actor, id); err != nil {
		log.Printf("[Handler.DeleteSession] Error from service Delete for session %s: %v", id.Hex(), err)
		if err.Error() == "sesión no encontrada para eliminar" || errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sesión no encontrada"})
		} else if err.Error() == "no tienes permiso para eliminar esta sesión" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la sesión"})
		}
		return
	}

	log.Printf("[Handler.DeleteSession] Session %s deleted successfully.", id.Hex())
	c.JSON(http.StatusOK, gin.H{"message": "Sesión eliminada correctamente"})
}

// GET /api/workout-sessions/search → buscar sesiones
func (h *WorkoutSessionHandler) SearchSessions(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}
	log.Printf("[Handler.SearchSessions] Request from actor %s (role %s)", actor.UserID.Hex(), actor.Role)

	filter := bson.M{}
	if routineIDStr := c.Query("routine_id"); routineIDStr != "" {
		if oid, err := primitive.ObjectIDFromHex(routineIDStr); err == nil {
			filter["routineId"] = oid
		} else {
			log.Printf("[Handler.SearchSessions] Invalid routine_id query param: %s", routineIDStr)

		}
	}

	opts := options.Find()

	log.Printf("[Handler.SearchSessions] Calling service Search with filter: %v", filter)
	resultsModels, err := h.sessionService.Search(actor, filter, opts)
	if err != nil {
		log.Printf("[Handler.SearchSessions] Error from service Search: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar sesiones"})
		return
	}

	resultsDTOs := make([]dto.WorkoutSessionResponse, 0, len(resultsModels))
	for _, model := range resultsModels {
		resultsDTOs = append(resultsDTOs, utils.ConvertWorkoutSessionModelToResponse(model))
	}

	log.Printf("[Handler.SearchSessions] Found %d sessions. Returning DTO list.", len(resultsDTOs))
	c.JSON(http.StatusOK, resultsDTOs)
}
