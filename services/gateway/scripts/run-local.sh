#!/bin/bash
set -euo pipefail

export APP_ENV=development
export LOG_LEVEL=debug
export GATEWAY_HTTP_PORT=8080
export REDIS_ADDR=localhost:6379
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
export JWKS_ENDPOINT=http://localhost:8081/.well-known/jwks.json

export UPSTREAM_AUTH_SERVICE=localhost:8081
export UPSTREAM_CATALOG_SERVICE=localhost:8082
export UPSTREAM_CART_SERVICE=localhost:8083
export UPSTREAM_ORDER_SERVICE=localhost:8084
export UPSTREAM_INVENTORY_SERVICE=localhost:8085
export UPSTREAM_PAYMENT_SERVICE=localhost:8086
export UPSTREAM_SEARCH_SERVICE=localhost:8087
export UPSTREAM_RECOMMENDATION_SERVICE=localhost:8088

echo "Starting Gateway in development mode..."
go run ./cmd/server
