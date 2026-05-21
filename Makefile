.PHONY: all build test lint proto docker run dev clean help

GO ?= go
GOLANGCI_LINT ?= golangci-lint
PROTOC ?= protoc
DOCKER ?= docker

GO_FLAGS ?= -ldflags="-w -s"
CGO_ENABLED ?= 0

# Detect modules from go.work
MODULES := $(shell grep '\./' go.work | sed 's/\t*\.\///' | tr -d '\n')

all: lint test build

help:
	@echo "Targets:"
	@echo "  build       - Build all modules"
	@echo "  test        - Run all unit tests"
	@echo "  lint        - Run golangci-lint"
	@echo "  proto       - Generate protobuf code"
	@echo "  docker      - Build all Docker images"
	@echo "  run         - Run docker-compose"
	@echo "  dev         - Start local dev environment"
	@echo "  clean       - Clean build artifacts"
	@echo "  tidy        - Run go mod tidy on all modules"
	@echo "  security    - Run security scanning"

build:
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build ./...

test:
	CGO_ENABLED=$(CGO_ENABLED) $(GO) test ./... -count=1 -timeout=120s

test-race:
	CGO_ENABLED=1 $(GO) test -race ./... -count=1 -timeout=180s

test-integration:
	CGO_ENABLED=$(CGO_ENABLED) $(GO) test ./... -tags=integration -count=1 -timeout=300s

lint:
ifdef GOLANGCI_LINT
	$(GOLANGCI_LINT) run ./... --timeout=5m
else
	@echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
endif

proto:
	@for dir in proto/*/; do \
		if [ -d "$$dir" ]; then \
			echo "Generating proto in $$dir..."; \
			$(PROTOC) --proto_path=proto --go_out=proto --go_opt=paths=source_relative \
				--go-grpc_out=proto --go-grpc_opt=paths=source_relative \
				$$dir/*.proto; \
		fi; \
	done

docker:
	@for service in services/*/; do \
		name=$$(basename $$service); \
		echo "Building $$name..."; \
		$(DOCKER) build -t ghcr.io/shopee-clone/$$name:latest -f $$service/Dockerfile $$service; \
	done
	@for platform in platforms/*/; do \
		name=$$(basename $$platform); \
		echo "Building $$name..."; \
		if [ -f "$$platform/deployments/Dockerfile" ]; then \
			$(DOCKER) build -t ghcr.io/shopee-clone/$$name:latest -f $$platform/deployments/Dockerfile $$platform; \
		elif [ -f "$$platform/Dockerfile" ]; then \
			$(DOCKER) build -t ghcr.io/shopee-clone/$$name:latest -f $$platform/Dockerfile $$platform; \
		fi; \
	done

run:
	$(DOCKER)-compose up -d

dev: run
	@echo "Dev environment started. Run individual services with: go run ./cmd/server"

stop:
	$(DOCKER)-compose down

clean:
	$(GO) clean -cache -testcache
	rm -rf bin/

tidy:
	@for mod in $(MODULES); do \
		echo "Tidying $$mod..."; \
		(cd $$mod && $(GO) mod tidy) || true; \
	done

security:
	@echo "Running security checks..."
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Install with: go install github.com/securego/gosec/v2/cmd/gosec@latest"; \
	fi

benchmark:
	CGO_ENABLED=$(CGO_ENABLED) $(GO) test -bench=. -benchmem ./...

.PHONY: build test test-race test-integration lint proto docker run dev stop clean tidy security benchmark
