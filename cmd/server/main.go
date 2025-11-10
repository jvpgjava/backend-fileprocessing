package main

import (
	"log"
	"os"

	_ "backend-fileprocessing/docs"
	"backend-fileprocessing/internal/config"
	"backend-fileprocessing/internal/server"

	"github.com/joho/godotenv"
)

// @title Backend File Processing API
// @version 1.0.0
// @description API para processamento de arquivos com extração de texto usando OCR e processamento nativo
// @host localhost:9091
// @BasePath /
// @schemes http
func main() {
	// Carregar variáveis de ambiente do arquivo .env (se existir)
	// Isso facilita desenvolvimento local - em produção use variáveis de ambiente reais
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️ Arquivo .env não encontrado (isso é normal em produção)")
	}

	// Carregar configurações
	cfg := config.Load()

	router := server.NewRouter(cfg)

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
