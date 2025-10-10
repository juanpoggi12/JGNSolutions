package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juanpoggi12/JGNSolutions/backend/services"
)

type AdminLogsHandler struct {
	logService *services.LogService
}

func NewAdminLogsHandler(logService *services.LogService) *AdminLogsHandler {
	return &AdminLogsHandler{logService: logService}
}

// Obtener todos los logs (solo para admins)
func (h *AdminLogsHandler) GetLogs(c *gin.Context) {
	logs, err := h.logService.ListarLogs(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}
