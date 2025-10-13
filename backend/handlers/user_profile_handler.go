package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/services"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserProfileHandler maneja las peticiones HTTP relacionadas con los perfiles de usuario
type UserProfileHandler struct {
	userProfileService *services.UserProfileService
}

// Constructor
func NewUserProfileHandler(userProfileService *services.UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{userProfileService: userProfileService}
}

// GET /api/profile → Obtener perfil del usuario autenticado
func (h *UserProfileHandler) GetProfile(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	perfil, err := h.userProfileService.GetProfile(actor)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, perfil)
}

func (h *UserProfileHandler) UpdateProfile(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	var req dto.UserProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	updated, err := h.userProfileService.UpdateProfile(actor, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// POST /api/profile → Crear perfil (si aún no existe)
func (h *UserProfileHandler) CreateProfile(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	// ✅ Bind al DTO
	var req dto.UserProfileCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	created, err := h.userProfileService.CreateProfile(actor, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// DELETE /api/profile/:id → Eliminar perfil (solo admin o el mismo usuario)
func (h *UserProfileHandler) DeleteProfile(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	id := c.Param("id")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Si no es admin y quiere borrar otro perfil, denegamos
	if actor.Role != "admin" && actor.UserID != objectID {
		c.JSON(http.StatusForbidden, gin.H{"error": "no tienes permiso para eliminar este perfil"})
		return
	}

	if err := h.userProfileService.DeleteProfile(actor, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Perfil eliminado correctamente"})
}
