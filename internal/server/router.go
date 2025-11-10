package server

import (
	"backend-fileprocessing/internal/config"
	"backend-fileprocessing/internal/handlers"
	"backend-fileprocessing/internal/middleware"
	"backend-fileprocessing/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter configura e retorna um *gin.Engine pronto para uso.
func NewRouter(cfg *config.Config) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	fileService := services.NewFileService()

	fileHandler := handlers.NewFileHandler(fileService)
	healthHandler := handlers.NewHealthHandler()

	setupRoutes(router, fileHandler, healthHandler)

	return router
}

func setupRoutes(router *gin.Engine, fileHandler *handlers.FileHandler, healthHandler *handlers.HealthHandler) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.HealthCheck)
		v1.GET("/status", healthHandler.Status)

		files := v1.Group("/files")
		{
			files.POST("/process", fileHandler.ProcessFile)
			files.GET("/supported-types", fileHandler.GetSupportedTypes)
		}
	}
}
