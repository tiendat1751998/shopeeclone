.PHONY: all build test lint clean docker

APP_NAME := shopee-payment
VERSION := $(shell git describe --tags --always 2>/dev/null || echo "dev")

all: lint test build
build:
	CGO_ENABLED=0 go build -ldflags="-w -s -X main.version=$(VERSION)" -o bin/$(APP_NAME) ./cmd/server
test:
	go test -v -race -coverprofile=coverage.out ./...
lint:
	golangci-lint run ./...
clean:
	rm -rf bin/ coverage.out
docker:
	docker build -t shopee-clone/$(APP_NAME):$(VERSION) .
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/payment/v1/payment.proto
run-local:
	go run ./cmd/server
fmt:
	go fmt ./...
