package middleware

import (
	"log"
	"net/http"

	"backend-fileprocessing/internal/models"

	"github.com/gin-gonic/gin"
)

// Recovery middleware de recovery
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.Printf("Panic recovered: %s", err)
		}

		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"INTERNAL_ERROR",
			"Erro interno do servidor",
			"Ocorreu um erro inesperado",
		))
		c.Abort()
	})
}
