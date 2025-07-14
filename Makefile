.PHONY: build run test clean docker-up docker-down migrate-up migrate-down deps lint fmt vet

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=meteo-api
BINARY_PATH=./bin/$(BINARY_NAME)
MAIN_PATH=./cmd/server

# Build the application
build:
	$(GOBUILD) -o $(BINARY_PATH) -v $(MAIN_PATH)

# Run the application
run:
	$(GOBUILD) -o $(BINARY_PATH) -v $(MAIN_PATH)
	$(BINARY_PATH)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_PATH)
	rm -f coverage.out coverage.html

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
fmt:
	$(GOCMD) fmt ./...

# Vet code
vet:
	$(GOCMD) vet ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Start services with Docker Compose
docker-up:
	docker-compose up -d

# Stop services
docker-down:
	docker-compose down

# Build and start services
docker-build:
	docker-compose up --build -d

# View logs
docker-logs:
	docker-compose logs -f

# Database migrations
migrate-up:
	migrate -path internal/database/migrations -database "postgres://postgres:postgres@localhost:5432/meteo_db?sslmode=disable" up

migrate-down:
	migrate -path internal/database/migrations -database "postgres://postgres:postgres@localhost:5432/meteo_db?sslmode=disable" down

# Create new migration
migrate-create:
	migrate create -ext sql -dir internal/database/migrations -seq $(name)

# Development setup
dev-setup: deps
	cp .env.example .env
	docker-compose up -d postgres
	sleep 5
	make migrate-up

# Full development workflow
dev: dev-setup
	make run

# Production build
build-prod:
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -a -installsuffix cgo -o $(BINARY_PATH) $(MAIN_PATH)

# Install tools
install-tools:
	$(GOGET) -u github.com/golang-migrate/migrate/v4/cmd/migrate
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Build and run the application"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Download dependencies"
	@echo "  fmt            - Format code"
	@echo "  vet            - Vet code"
	@echo "  lint           - Lint code"
	@echo "  docker-up      - Start services with Docker"
	@echo "  docker-down    - Stop Docker services"
	@echo "  docker-build   - Build and start services"
	@echo "  docker-logs    - View Docker logs"
	@echo "  migrate-up     - Run database migrations"
	@echo "  migrate-down   - Rollback database migrations"
	@echo "  migrate-create - Create new migration (use: make migrate-create name=migration_name)"
	@echo "  dev-setup      - Setup development environment"
	@echo "  dev            - Full development workflow"
	@echo "  build-prod     - Production build"
	@echo "  install-tools  - Install development tools"
	@echo "  help           - Show this help message"