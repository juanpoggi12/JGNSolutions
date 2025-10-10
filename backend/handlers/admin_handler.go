package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juanpoggi12/JGNSolutions/backend/services"
)

// AdminHandler maneja las rutas de estadísticas y administración
type AdminHandler struct {
	adminService *services.AdminService
}

// Constructor
func NewAdminHandler(adminService *services.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// GET /api/admin/users/count → Total de usuarios
func (h *AdminHandler) CountUsers(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	count, err := h.adminService.CountUsers(actor)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_users": count})
}

// GET /api/admin/exercises/count → Total de ejercicios
func (h *AdminHandler) CountExercises(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	count, err := h.adminService.CountExercises(actor)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_exercises": count})
}

// GET /api/admin/routines/count → Total de rutinas
func (h *AdminHandler) CountRoutines(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	count, err := h.adminService.CountRoutines(actor)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_routines": count})
}

// GET /api/admin/workouts/count → Total de sesiones de entrenamiento
func (h *AdminHandler) CountWorkoutSessions(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	count, err := h.adminService.CountWorkoutSessions(actor)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_workout_sessions": count})
}

// GET /api/admin/users → Lista básica de usuarios (solo admin)
func (h *AdminHandler) ListUsers(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	users, err := h.adminService.ListUsers(actor)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GET /api/admin/exercises/top?limit=5 → Ejercicios más usados
func (h *AdminHandler) TopExercises(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	stats, err := h.adminService.TopExercises(actor, limit)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GET /api/admin/routines/top?limit=5 → Rutinas más usadas
func (h *AdminHandler) TopRoutines(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	stats, err := h.adminService.TopRoutines(actor, limit)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
