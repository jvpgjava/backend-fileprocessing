# Makefile para Backend File Processing

.PHONY: build run test clean deps lint format

# Variáveis
BINARY_NAME=backend-fileprocessing
BUILD_DIR=build
MAIN_PATH=cmd/server/main.go

# Build
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# Run
run:
	@echo "🚀 Running $(BINARY_NAME)..."
	@go run $(MAIN_PATH)

# Run with hot reload (requer air)
dev:
	@echo "🔥 Running with hot reload..."
	@air

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	@go mod tidy
	@go mod download

# Test
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Test with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Lint
lint:
	@echo "🔍 Running linter..."
	@golangci-lint run

# Format
format:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@goimports -w .

# Clean
clean:
	@echo "🧹 Cleaning build files..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# Install tools
install-tools:
	@echo "🛠️ Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger docs
swagger:
	@echo "📚 Generating Swagger documentation..."
	@swag init -g cmd/server/main.go -o docs

# Local development only

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the binary"
	@echo "  run           - Run the application"
	@echo "  dev           - Run with hot reload (requires air)"
	@echo "  deps          - Install dependencies"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  lint          - Run linter"
	@echo "  format        - Format code"
	@echo "  clean         - Clean build files"
	@echo "  install-tools - Install development tools"
	@echo "  swagger       - Generate Swagger documentation"
	@echo "  help          - Show this help"
