package main

import (
	"log"
	"os"

	"backend-fileprocessing/docs"
	"backend-fileprocessing/internal/config"
	"backend-fileprocessing/internal/handlers"
	"backend-fileprocessing/internal/middleware"
	"backend-fileprocessing/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Backend File Processing API
// @version 1.0.0
// @description API para processamento de arquivos com extração de texto usando OCR e processamento nativo
// @host localhost:9091
// @BasePath /
// @schemes http
func main() {
	// Carregar configurações
	cfg := config.Load()

	// Configurar Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Criar router
	router := gin.New()

	// Middleware global
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// Inicializar serviços
	fileService := services.NewFileService()

	// Inicializar handlers
	fileHandler := handlers.NewFileHandler(fileService)
	healthHandler := handlers.NewHealthHandler()

	// Configurar rotas
	setupRoutes(router, fileHandler, healthHandler)

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}

	log.Printf("🚀 Servidor iniciado na porta %s", port)
	log.Printf("📁 Serviço de processamento de arquivos ativo")
	log.Printf("🔗 Health check: http://localhost:%s/api/v1/health", port)
	log.Printf("📤 Processar arquivo: http://localhost:%s/api/v1/files/process", port)
	log.Printf("📚 Documentação Swagger: http://localhost:%s/swagger/index.html", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}

func setupRoutes(router *gin.Engine, fileHandler *handlers.FileHandler, healthHandler *handlers.HealthHandler) {
	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", healthHandler.HealthCheck)
		v1.GET("/status", healthHandler.Status)

		// File processing
		files := v1.Group("/files")
		{
			files.POST("/process", fileHandler.ProcessFile)
			files.GET("/supported-types", fileHandler.GetSupportedTypes)
		}
	}
}
