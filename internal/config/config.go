package config

import (
	"os"
)

// Config estrutura de configuração
type Config struct {
	Port        string
	Environment string
	LogLevel    string
	MaxFileSize int64
}

// Load carrega configurações do ambiente
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "9091"),
		Environment: getEnv("GIN_MODE", "debug"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		MaxFileSize: 25 * 1024 * 1024, // 25MB (aumentado)
	}
}

// getEnv obtém variável de ambiente com valor padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
