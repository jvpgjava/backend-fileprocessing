package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"backend-fileprocessing/internal/models"

	"github.com/gin-gonic/gin"
)

// Recovery middleware de recovery
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Log detalhado do panic
		log.Printf("❌ PANIC RECOVERED: %v", recovered)
		log.Printf("❌ Stack trace:\n%s", debug.Stack())
		
		var errMsg string
		if err, ok := recovered.(string); ok {
			errMsg = err
		} else if err, ok := recovered.(error); ok {
			errMsg = err.Error()
		} else {
			errMsg = fmt.Sprintf("%v", recovered)
		}

		log.Printf("❌ Erro: %s", errMsg)
		log.Printf("❌ Path: %s %s", c.Request.Method, c.Request.URL.Path)

		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"INTERNAL_ERROR",
			fmt.Sprintf("Erro interno do servidor: %s", errMsg),
			"Verifique os logs do servidor para mais detalhes",
		))
		c.Abort()
	})
}
