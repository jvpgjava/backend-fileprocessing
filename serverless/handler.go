package serverless

import (
	"net/http"
	"sync"

	"backend-fileprocessing/internal/config"
	"backend-fileprocessing/internal/server"

	"github.com/joho/godotenv"
)

var (
	once    sync.Once
	handler http.Handler
)

func initHandler() {
	_ = godotenv.Load()
	cfg := config.Load()
	handler = server.NewRouter(cfg)
}

// NewHandler retorna um http.Handler inicializado para uso em ambientes serverless.
func NewHandler() http.Handler {
	once.Do(initHandler)
	if handler == nil {
		initHandler()
	}
	return handler
}
