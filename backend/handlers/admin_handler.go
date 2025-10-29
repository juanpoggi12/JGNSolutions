package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juanpoggi12/JGNSolutions/backend/dto"
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

// GET /api/admin/routines/count → Total de rutinas

// GET /api/admin/workouts/count → Total de sesiones de entrenamiento

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

// --- 📜 LISTAR LOGS DEL SISTEMA (solo admin) ---
func (h *AdminHandler) ListLogs(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	logs, err := h.adminService.ListLogs(actor)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	id := c.Param("id")

	var req dto.AdminRoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	updated, err := h.adminService.UpdateUserRole(actor, id, req.Role)
	if err != nil {
		status := http.StatusForbidden
		if err.Error() == "rol inválido" || err.Error() == "usuario no encontrado" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *AdminHandler) MetricsSummary(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	summary, err := h.adminService.GetMetricsSummary(actor)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
