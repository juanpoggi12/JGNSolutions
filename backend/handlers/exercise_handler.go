package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/services"
)

// ExerciseHandler maneja las peticiones HTTP relacionadas con los ejercicios
type ExerciseHandler struct {
	exerciseService *services.ExerciseService
}

// Constructor
func NewExerciseHandler(exerciseService *services.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{exerciseService: exerciseService}
}

func (h *ExerciseHandler) ListExercises(c *gin.Context) {
	actor := services.Actor{
		UserID: parseObjectID(c.GetString("userId")),
		Role:   c.GetString("role"),
	}

	var query dto.ExerciseCatalogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros inválidos"})
		return
	}

	catalog, err := h.exerciseService.ListExercises(actor, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, catalog)
}

// POST /api/exercises → Crear nuevo ejercicio (solo admin)
func (h *ExerciseHandler) CreateExercise(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	var req dto.ExerciseCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	exercise, err := h.exerciseService.CreateExercise(actor, req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, exercise)
}

// GET /api/exercises/:id → Obtener ejercicio por ID
func (h *ExerciseHandler) GetExerciseByID(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	id := c.Param("id")
	exercise, err := h.exerciseService.GetExerciseByID(actor, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exercise)
}

// PUT /api/exercises/:id → Actualizar ejercicio (solo admin)
func (h *ExerciseHandler) UpdateExercise(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	id := c.Param("id")
	var req dto.ExerciseUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	updated, err := h.exerciseService.UpdateExercise(actor, id, req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DELETE /api/exercises/:id → Eliminar ejercicio (solo admin)
func (h *ExerciseHandler) DeleteExercise(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	id := c.Param("id")
	if err := h.exerciseService.DeleteExercise(actor, id); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ejercicio eliminado correctamente"})
}

// GET /api/exercises/search → Buscar ejercicios (admin y user)
func (h *ExerciseHandler) SearchExercises(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	var search dto.ExerciseSearchRequest
	if err := c.ShouldBindQuery(&search); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros de búsqueda inválidos"})
		return
	}

	exercises, err := h.exerciseService.SearchExercises(actor, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exercises)
}
