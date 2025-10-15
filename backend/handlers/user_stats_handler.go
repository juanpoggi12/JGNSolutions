package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juanpoggi12/JGNSolutions/backend/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStatsHandler struct {
	statsService services.UserStatsServiceInterface
}

func NewUserStatsHandler(statsService services.UserStatsServiceInterface) *UserStatsHandler {
	return &UserStatsHandler{statsService: statsService}
}

func (h *UserStatsHandler) GetWorkoutSummary(c *gin.Context) {
	actor := services.Actor{
		UserID: parseObjectID(c.GetString("userId")),
		Role:   c.GetString("role"),
	}

	from, err := parseQueryDate(c.Query("from"), false)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido en 'from'"})
		return
	}
	to, err := parseQueryDate(c.Query("to"), true)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido en 'to'"})
		return
	}

	summary, err := h.statsService.GetWorkoutSummary(actor, from, to)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (h *UserStatsHandler) GetFrequency(c *gin.Context) {
	actor := services.Actor{
		UserID: parseObjectID(c.GetString("userId")),
		Role:   c.GetString("role"),
	}

	from, err := parseQueryDate(c.Query("from"), false)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido en 'from'"})
		return
	}
	to, err := parseQueryDate(c.Query("to"), true)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido en 'to'"})
		return
	}

	granularity := c.DefaultQuery("granularity", "week")

	resp, err := h.statsService.GetFrequency(actor, from, to, granularity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserStatsHandler) GetTopRoutines(c *gin.Context) {
	actor := services.Actor{
		UserID: parseObjectID(c.GetString("userId")),
		Role:   c.GetString("role"),
	}

	limitVal := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitVal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetro 'limit' inválido"})
		return
	}

	routines, err := h.statsService.GetTopRoutines(actor, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, routines)
}

func (h *UserStatsHandler) GetProgress(c *gin.Context) {
	actor := services.Actor{
		UserID: parseObjectID(c.GetString("userId")),
		Role:   c.GetString("role"),
	}

	exerciseID := c.Query("exerciseId")
	if exerciseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "exerciseId es obligatorio"})
		return
	}
	exID, err := primitive.ObjectIDFromHex(exerciseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "exerciseId inválido"})
		return
	}

	from, err := parseQueryDate(c.Query("from"), false)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido en 'from'"})
		return
	}
	to, err := parseQueryDate(c.Query("to"), true)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido en 'to'"})
		return
	}

	progress, err := h.statsService.GetExerciseProgress(actor, exID, from, to)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, progress)
}

func parseQueryDate(value string, end bool) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}
	if end {
		parsed = parsed.Add(24 * time.Hour)
	}
	return &parsed, nil
}
