.PHONY: all build test lint clean docker proto migrate

APP_NAME := shopee-order
VERSION := $(shell git describe --tags --always 2>/dev/null || echo "dev")
DOCKER_REGISTRY := shopee-clone
GOFLAGS := -ldflags="-w -s -X main.version=$(VERSION)"

all: lint test build

build:
	CGO_ENABLED=0 go build $(GOFLAGS) -o bin/$(APP_NAME) ./cmd/server

test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

test-unit:
	go test -v -race ./internal/... ./domain/...

test-integration:
	go test -v -race ./tests/integration/...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ coverage.out

docker:
	docker build -t $(DOCKER_REGISTRY)/$(APP_NAME):$(VERSION) .

docker-push:
	docker push $(DOCKER_REGISTRY)/$(APP_NAME):$(VERSION)

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	    proto/order/v1/order.proto

migrate-up:
	migrate -path migrations -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)" up

migrate-down:
	migrate -path migrations -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)" down

run-local:
	go run ./cmd/server

fmt:
	go fmt ./...
	goimports -w .

vet:
	go vet ./...

generate:
	go generate ./...
