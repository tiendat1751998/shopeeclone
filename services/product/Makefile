.PHONY: build test lint proto docker run dev migrate

# Variables
APP_NAME := product-service
DOCKER_REGISTRY := ghcr.io/shopee-clone
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -w -s"

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOVET := $(GOCMD) vet
GOMOD := $(GOCMD) mod
BINARY_NAME := server
MAIN_PATH := ./cmd/server

# Build the binary
build:
	@echo "Building $(APP_NAME)..."
	CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: bin/$(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -func=coverage.out

# Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, running go vet instead..."; \
		$(GOVET) ./...; \
	fi

# Generate protobuf code
proto:
	@echo "Generating protobuf code..."
	@if command -v protoc > /dev/null; then \
		protoc --go_out=. --go_opt=module=github.com/shopee-clone/shopee/services/product \
			--go-grpc_out=. --go-grpc_opt=module=github.com/shopee-clone/shopee/services/product \
			api/proto/*.proto; \
		echo "Protobuf generation complete."; \
	else \
		echo "protoc not found. Install from https://grpc.io/docs/protoc-installation/"; \
		exit 1; \
	fi

# Build Docker image
docker:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_REGISTRY)/$(APP_NAME):$(VERSION) .
	docker tag $(DOCKER_REGISTRY)/$(APP_NAME):$(VERSION) $(DOCKER_REGISTRY)/$(APP_NAME):latest
	@echo "Docker image built: $(DOCKER_REGISTRY)/$(APP_NAME):$(VERSION)"

# Run the service locally
run: build
	@echo "Running $(APP_NAME)..."
	./bin/$(BINARY_NAME)

# Run in development mode with hot reload
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Falling back to run..."; \
		$(MAKE) run; \
	fi

# Database migrations
migrate:
	@echo "Running database migrations..."
	@if command -v migrate > /dev/null; then \
		migrate -path migrations -database "$(DATABASE_URL)" up; \
	else \
		echo "migrate tool not found. Install with: go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; \
		exit 1; \
	fi

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/ coverage.out
	$(GOCMD) clean

# Show help
help:
	@echo "Available targets:"
	@echo "  build   - Build the binary"
	@echo "  test    - Run tests with coverage"
	@echo "  lint    - Run linter (golangci-lint or go vet)"
	@echo "  proto   - Generate protobuf code"
	@echo "  docker  - Build Docker image"
	@echo "  run     - Build and run locally"
	@echo "  dev     - Run with hot reload (air)"
	@echo "  migrate - Run database migrations"
	@echo "  deps    - Download and tidy dependencies"
	@echo "  clean   - Clean build artifacts"
	@echo "  help    - Show this help"
