package handlers

import (
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juanpoggi12/JGNSolutions/backend/dto"
	"github.com/juanpoggi12/JGNSolutions/backend/models"
	"github.com/juanpoggi12/JGNSolutions/backend/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserProfileHandler maneja las peticiones HTTP relacionadas con los perfiles de usuario
type UserProfileHandler struct {
	userProfileService *services.UserProfileService
	userService        *services.UserService
}

func NewUserProfileHandler(userProfileService *services.UserProfileService, userService *services.UserService) *UserProfileHandler {
	return &UserProfileHandler{userProfileService: userProfileService, userService: userService}
}

// GET /api/profile → Obtener perfil del usuario autenticado con DTO amigable
func (h *UserProfileHandler) GetMyProfile(c *gin.Context) {
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

	c.JSON(http.StatusOK, toProfileResponse(perfil))
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

// PUT /api/profile/:id? → Actualizar perfil propio o de otro usuario (si es admin)
func (h *UserProfileHandler) UpdateProfile(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")
	paramID := c.Param("id") // opcional: si no viene, se usa el propio

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	// Bind del cuerpo JSON
	var req dto.UserProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Si no se pasó :id, actualiza su propio perfil
	targetID := userID
	if paramID != "" {
		targetID = paramID
	}

	updated, err := h.userProfileService.UpdateProfile(actor, targetID, req)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "no encontrado") {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "permiso") {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.UserProfileResponse{
		ID:        updated.ID.Hex(),
		FullName:  updated.FullName,
		BirthDate: updated.BirthDate.Format("2006-01-02"),
		WeightKg:  updated.WeightKg,
		HeightCm:  updated.HeightCm,
		Level:     string(updated.Level),
		Goal:      string(updated.Goal),
		UpdatedAt: updated.UpdatedAt.Format(time.RFC3339),
	})
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

// POST /api/profile/change-password → cambiar contraseña del usuario autenticado
func (h *UserProfileHandler) ChangePassword(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetString("userId")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	if actor.UserID.IsZero() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	var req dto.ProfileChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	changeReq := dto.ChangePasswordRequest{
		OldPassword: req.CurrentPassword,
		NewPassword: req.NewPassword,
	}

	if err := h.userService.ChangePassword(actor, changeReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func mapProfileUpdateRequest(req dto.ProfileUpdateRequest) dto.UserProfileUpdateRequest {
	var out dto.UserProfileUpdateRequest

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name != "" {
			out.FullName = &name
		}
	}
	if req.BirthDate != nil {
		date := strings.TrimSpace(*req.BirthDate)
		if date != "" {
			out.BirthDate = &date
		}
	}
	if req.Weight != nil {
		weight := *req.Weight
		out.WeightKg = &weight
	}
	if req.Height != nil {
		h := int(math.Round(*req.Height))
		out.HeightCm = &h
	}
	if req.Level != nil {
		level := normalizeLevel(*req.Level)
		out.Level = &level
	}
	if req.Goal != nil {
		goal := normalizeGoal(*req.Goal)
		if goal != "" {
			out.Goal = &goal
		}
	}

	return out
}

func normalizeLevel(level string) string {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "beginner":
		return "PRINCIPIANTE"
	case "intermediate":
		return "INTERMEDIO"
	case "advanced":
		return "AVANZADO"
	default:
		return strings.ToUpper(level)
	}
}

func toProfileResponse(profile models.UserProfile) dto.ProfileResponse {
	birthDate := profile.BirthDate.Format("2006-01-02")

	weight := profile.WeightKg
	height := float64(profile.HeightCm)

	level := mapLevelForResponse(string(profile.Level))
	goal := mapGoalForResponse(string(profile.Goal))

	return dto.ProfileResponse{
		ID:        profile.ID.Hex(),
		Name:      profile.FullName,
		BirthDate: birthDate,
		Weight:    weight,
		Height:    height,
		Level:     level,
		Goal:      goal,
		UpdatedAt: profile.UpdatedAt.Format(time.RFC3339),
	}
}

func mapLevelForResponse(level string) string {
	switch strings.ToUpper(level) {
	case "PRINCIPIANTE":
		return "beginner"
	case "INTERMEDIO":
		return "intermediate"
	case "AVANZADO":
		return "advanced"
	default:
		return strings.ToLower(level)
	}
}

func normalizeGoal(goal string) string {
	switch strings.ToLower(strings.TrimSpace(goal)) {
	case "lose_weight", "perder_peso", "perderpeso", "weight_loss":
		return "PERDER_PESO"
	case "gain_muscle", "ganar_musculo", "ganarmusculo", "muscle_gain":
		return "GANAR_MUSCULO"
	case "maintain", "mantenerse", "mantener", "maintenance":
		return "MANTENERSE"
	default:
		return strings.ToUpper(strings.TrimSpace(goal))
	}
}

func mapGoalForResponse(goal string) string {
	switch strings.ToUpper(goal) {
	case "PERDER_PESO":
		return "lose_weight"
	case "GANAR_MUSCULO":
		return "gain_muscle"
	case "MANTENERSE":
		return "maintain"
	default:
		return strings.ToLower(goal)
	}
}
