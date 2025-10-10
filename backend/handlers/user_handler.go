package handlers

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/services"
	"github.com/juanpoggi12/JGNSolutions/backend/dto"

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
func (h *UserHandler) CreateUser(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	created, err := h.userService.CreateUser(actor, user)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

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
func (h *UserHandler) UpdateUser(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	id := c.Param("id")
	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	updated, err := h.userService.UpdateUser(actor, id, req)

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

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

// Cambiar contraseña del usuario autenticado
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest

	// 1️⃣ Validar el JSON recibido
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2️⃣ Obtener el ID del usuario del contexto (lo setea el middleware JWT)
	userIDHex, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	// 3️⃣ Convertir el ID de string a ObjectID
	userID, err := primitive.ObjectIDFromHex(userIDHex.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	// 4️⃣ Llamar al servicio para cambiar la contraseña
	if err := h.userService.ChangePassword(c, userID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 5️⃣ Responder al cliente
	c.JSON(http.StatusOK, gin.H{
		"message": "Contraseña actualizada correctamente",
	})
}