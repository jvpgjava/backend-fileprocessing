package handler

import (
	"net/http"
	"sync"

	"backend-fileprocessing/internal/config"
	"backend-fileprocessing/internal/server"

	"github.com/joho/godotenv"
)

var (
	routerOnce sync.Once
	router     http.Handler
)

func initRouter() {
	_ = godotenv.Load()
	cfg := config.Load()
	router = server.NewRouter(cfg)
}

// Handler é o entrypoint requisitado pela Vercel.
func Handler(w http.ResponseWriter, r *http.Request) {
	routerOnce.Do(initRouter)
	router.ServeHTTP(w, r)
}
