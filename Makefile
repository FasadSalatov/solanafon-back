.PHONY: help install run seed dev build clean test

help: ## Show this help
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies
	go mod download
	go mod verify

run: ## Run the server
	go run cmd/server/main.go

seed: ## Seed the database with initial data
	go run cmd/seed/main.go

dev: ## Run in development mode with hot reload (requires air)
	air

build: ## Build the application
	go build -o bin/server cmd/server/main.go

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

test: ## Run tests
	go test -v ./...

lint: ## Run linter
	golangci-lint run

docker-up: ## Start Docker services
	docker-compose up -d

docker-down: ## Stop Docker services
	docker-compose down

docker-logs: ## Show Docker logs
	docker-compose logs -f
