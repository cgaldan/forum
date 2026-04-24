.PHONY: help build run test clean docker-build docker-run docker-compose-up docker-compose-down migrate

help:
	@echo "Available commands:"
	@echo "  make build                - Build the application"
	@echo "  make run                  - Run the application"
	@echo "  make test                 - Run tests"
	@echo "  make clean                - Clean build artifacts"
	@echo "  make docker-build         - Build Docker image (needs repo root as context)"
	@echo "  make docker-run           - Run container (port 8000)"
	@echo "  make docker-compose-up    - docker compose up --build (from repo root)"
	@echo "  make docker-compose-down  - docker compose down (from repo root)"
	@echo "  make deps                 - Install dependencies"

build:
	@echo "Building..."
	@go build -o bin/forum-backend ./cmd/server

run:
	@echo "Running..."
	@go run ./cmd/server/main.go

test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

docker-build:
	@echo "Building Docker image..."
	@docker build -f Dockerfile -t forum-backend:latest ..

docker-run:
	@echo "Running Docker container..."
	@docker run --rm -p 8000:8000 --name forum-backend forum-backend:latest

docker-compose-up:
	@docker compose -f ../docker-compose.yml up --build

docker-compose-down:
	@docker compose -f ../docker-compose.yml down

deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy