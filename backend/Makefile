.PHONY: help build run test clean docker-build docker-run migrate

# Default target
help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run Docker container"
	@echo "  make migrate      - Run database migrations"
	@echo "  make lint         - Run linter"
	@echo "  make fmt          - Format code"

# Build the application
build:
	@echo "Building..."
	@go build -o bin/forum-backend ./cmd/server

# Run the application
run:
	@echo "Running..."
	@go run ./cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t forum-backend:latest .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8000:8000 --name forum-backend forum-backend:latest

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Development mode with auto-reload
dev:
	@echo "Running in development mode..."
	@air

# Generate API documentation
docs:
	@echo "Generating API documentation..."
	@swag init -g cmd/server/main.go

