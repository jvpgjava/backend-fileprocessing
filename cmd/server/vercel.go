package main

import (
	"net/http"
	"sync"

	"backend-fileprocessing/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	vercelOnce   sync.Once
	vercelRouter http.Handler
)

func initVercelRouter() {
	// Tenta carregar .env (ignora erro em produção)
	_ = godotenv.Load()

	cfg := config.Load()
	vercelRouter = newRouter(cfg)
}

// Handler é o ponto de entrada exportado para a Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	vercelOnce.Do(initVercelRouter)

	if vercelRouter == nil {
		// Fallback defensivo
		cfg := config.Load()
		vercelRouter = newRouter(cfg)
	}

	vercelRouter.ServeHTTP(w, r)
}
