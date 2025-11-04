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
	paramID := c.Param("id")

	actor := services.Actor{
		UserID: parseObjectID(userID),
		Role:   role,
	}

	var frontendReq dto.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&frontendReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	serviceReq := mapProfileUpdateRequest(frontendReq)

	targetID := userID
	if paramID != "" {
		targetID = paramID
	}

	updated, err := h.userProfileService.UpdateProfile(actor, targetID, serviceReq)
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
