package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juanpoggi12/JGNSolutions/backend/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserHandler maneja las peticiones HTTP relacionadas con usuarios
type UserHandler struct {
	userService *services.UserService
}

// Constructor
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// POST /api/users → Crear usuario (solo admin)

// GET /api/users/:id → Obtener usuario (admin o el mismo usuario)
func (h *UserHandler) GetUserByID(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	id := c.Param("id")
	user, err := h.userService.GetUserByID(actor, id)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// PUT /api/users/:id → Actualizar usuario (admin o el mismo usuario)

// DELETE /api/users/:id → Eliminar usuario (solo admin)
func (h *UserHandler) DeleteUser(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	id := c.Param("id")
	if err := h.userService.DeleteUser(actor, id); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado correctamente"})
}

// Helper para convertir string a ObjectID sin repetir código
func parseObjectID(id string) primitive.ObjectID {
	oid, _ := primitive.ObjectIDFromHex(id)
	return oid
}
