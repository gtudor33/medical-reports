.PHONY: help build run test clean docker-up docker-down migrate

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building medical-reports-api..."
	@go build -o medical-reports-api cmd/api/main.go
	@echo "Build complete!"

run: ## Run the application locally
	@echo "Running medical-reports-api..."
	@go run cmd/api/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f medical-reports-api
	@go clean

docker-up: ## Start all services with docker-compose
	@echo "Starting services..."
	@docker compose up -d

docker-down: ## Stop all services
	@echo "Stopping services..."
	@docker compose down

docker-logs: ## Show logs from all services
	@docker compose logs -f

docker-build: ## Build docker image
	@docker compose build

tidy: ## Tidy go modules
	@go mod tidy

fmt: ## Format code
	@go fmt ./...

vet: ## Run go vet
	@go vet ./...

lint: fmt vet ## Run formatters and linters

deps: ## Download dependencies
	@go mod download

dev: docker-up ## Start development environment
	@echo "Development environment ready!"
	@echo "API: http://localhost:8080"
	@echo "PostgreSQL: localhost:5432"

check: ## Check if code compiles
	@echo "Checking if code compiles..."
	@go build -o /dev/null cmd/api/main.go
	@echo "✓ Code compiles successfully!"
