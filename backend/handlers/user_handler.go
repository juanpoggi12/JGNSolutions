package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juanpoggi12/JGNSolutions/backend/dto"
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
func (h *UserHandler) CreateUser(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	// 1️⃣ Bindeamos al DTO, no al modelo
	var req dto.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// 2️⃣ Llamar al servicio (él se encarga del resto)
	created, err := h.userService.CreateUser(actor, req)
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
	// 1️⃣ Validar el JSON recibido
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// 2️⃣ Obtener los datos del usuario autenticado (set por middleware JWT)
	role := c.GetString("role")
	userID := c.GetString("userId") // ⚠️ corregido: antes usabas "userID"

	// 3️⃣ Crear el actor
	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	// 4️⃣ Llamar al servicio para cambiar la contraseña
	if err := h.userService.ChangePassword(actor, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 5️⃣ Responder al cliente
	c.JSON(http.StatusOK, gin.H{
		"message": "Contraseña actualizada correctamente",
	})
}
