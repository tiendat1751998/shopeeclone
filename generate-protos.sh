#!/bin/sh
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "${SCRIPT_DIR}"

echo "Starting protobuf compilation inside a golang Docker container..."

docker run --rm \
  -v "$(pwd)":/workspace \
  -w /workspace \
  golang:1.24-alpine \
  sh -c "
    apk add --no-cache protobuf-dev protoc git && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    
    echo 'Generating central catalog proto...' && \
    mkdir -p proto/catalog/v1 && \
    protoc --proto_path=proto/shopee \
           --go_out=proto \
           --go_opt=module=github.com/shopee-clone/shopee/proto \
           --go-grpc_out=proto \
           --go-grpc_opt=module=github.com/shopee-clone/shopee/proto \
           proto/shopee/catalog/v1/catalog.proto && \
           
    echo 'Generating services/order proto...' && \
    protoc --proto_path=services/order \
           --go_out=services/order \
           --go_opt=paths=source_relative \
           --go-grpc_out=services/order \
           --go-grpc_opt=paths=source_relative \
           services/order/proto/order/v1/order.proto && \
           
    echo 'Generating services/payment proto...' && \
    protoc --proto_path=services/payment \
           --go_out=services/payment \
           --go_opt=paths=source_relative \
           --go-grpc_out=services/payment \
           --go-grpc_opt=paths=source_relative \
           services/payment/proto/payment/v1/payment.proto && \
           
    echo 'Generating services/inventory proto...' && \
    protoc --proto_path=services/inventory \
           --go_out=services/inventory \
           --go_opt=paths=source_relative \
           --go-grpc_out=services/inventory \
           --go-grpc_opt=paths=source_relative \
           services/inventory/proto/inventory/v1/inventory.proto && \
           
    echo 'Generating services/product-catalog proto...' && \
    protoc --proto_path=services/product-catalog \
           --go_out=services/product-catalog \
           --go_opt=paths=source_relative \
           --go-grpc_out=services/product-catalog \
           --go-grpc_opt=paths=source_relative \
           services/product-catalog/proto/productcatalog/v1/catalog.proto && \
           
    echo 'Generating services/shipment proto...' && \
    protoc --proto_path=services/shipment \
           --go_out=services/shipment \
           --go_opt=paths=source_relative \
           --go-grpc_out=services/shipment \
           --go-grpc_opt=paths=source_relative \
           services/shipment/proto/shipment/v1/shipment.proto
  "

echo "All protobuf compilation completed successfully!"
