package handlers

import (
    "net/http"
    "time"

    "backend-fileprocessing/internal/models"

    "github.com/gin-gonic/gin"
)

// HealthHandler handler para health checks
type HealthHandler struct {
	startTime time.Time
}

// NewHealthHandler cria novo handler de health
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
	}
}

// HealthCheck verifica saúde do serviço
// @Summary Health check
// @Description Verifica a saúde do serviço
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	uptime := time.Since(h.startTime)
	
	response := models.HealthResponse{
		Status:    "ok",
		Service:   "backend-fileprocessing",
		Version:   "1.0.0",
		Timestamp: time.Now(),
		Uptime:    uptime.String(),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// Status retorna status detalhado do serviço
// @Summary Status detalhado
// @Description Retorna status detalhado do serviço
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/status [get]
func (h *HealthHandler) Status(c *gin.Context) {
	response := models.StatusResponse{
		Service:     "backend-fileprocessing",
		Version:     "1.0.0",
		Status:      "running",
		Timestamp:   time.Now(),
		Environment: gin.Mode(),
		Features: []string{
			"PDF Processing",
			"Image OCR",
			"Text Extraction",
			"DOCX Support",
			"REST API",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}
