package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/services"
)

// RoutineHandler maneja las operaciones sobre rutinas y sus ejercicios
type RoutineHandler struct {
	routineService *services.RoutineService
}

// Constructor
func NewRoutineHandler(routineService *services.RoutineService) *RoutineHandler {
	return &RoutineHandler{routineService: routineService}
}

// POST /api/routines → crear rutina
func (h *RoutineHandler) CreateRoutine(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	var req dto.RoutineCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	routine, err := h.routineService.CreateRoutine(actor, req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, routine)
}

// DELETE /api/routines/:id → eliminar rutina
func (h *RoutineHandler) DeleteRoutine(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	id := c.Param("id")
	if err := h.routineService.DeleteRoutine(actor, id); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rutina eliminada correctamente"})
}

// GET /api/routines/search → buscar rutinas
func (h *RoutineHandler) SearchRoutines(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	var search dto.RoutineSearchRequest
	if err := c.ShouldBindQuery(&search); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros inválidos"})
		return
	}

	results, err := h.routineService.SearchRoutines(actor, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// POST /api/routines/:id/exercises → agregar ejercicio a una rutina
func (h *RoutineHandler) AddExerciseToRoutine(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	var req dto.RoutineExerciseCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	added, err := h.routineService.AddExerciseToRoutine(actor, req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, added)
}

// PUT /api/routines/exercises/:id → actualizar ejercicio en una rutina
func (h *RoutineHandler) UpdateRoutineExercise(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	id := c.Param("id")
	var req dto.RoutineExerciseUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	updated, err := h.routineService.UpdateRoutineExercise(actor, id, req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DELETE /api/routines/exercises/:id → eliminar ejercicio de una rutina
func (h *RoutineHandler) DeleteRoutineExercise(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	id := c.Param("id")
	if err := h.routineService.DeleteRoutineExercise(actor, id); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ejercicio eliminado de la rutina"})
}

// GET /api/routines/exercises/search → listar ejercicios de una rutina
func (h *RoutineHandler) ListRoutineExercises(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	var search dto.RoutineExerciseSearchRequest
	if err := c.ShouldBindQuery(&search); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros inválidos"})
		return
	}

	results, err := h.routineService.ListRoutineExercises(actor, search)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// POST /api/routines/:id/duplicate → duplicar rutina
func (h *RoutineHandler) DuplicateRoutine(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	id := c.Param("id")

	var body struct {
		NewName string `json:"new_name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		body.NewName = ""
	}

	duplicated, err := h.routineService.DuplicateRoutine(actor, id, body.NewName)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, duplicated)
}
