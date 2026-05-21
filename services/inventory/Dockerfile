FROM golang:1.26.3-bookworm AS builder

RUN apt-get update && apt-get install -y git ca-certificates tzdata

WORKDIR /app/services/inventory

COPY packages/go-shared /app/packages/go-shared
COPY services/inventory/go.mod services/inventory/go.sum ./
RUN go mod download

COPY services/inventory/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o /app/bin/inventory \
    ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app

COPY --from=builder /app/bin/inventory /app/inventory
COPY --from=builder /app/services/inventory/migrations /app/migrations
COPY --from=builder /app/services/inventory/configs /app/configs

EXPOSE 8086 9096

ENTRYPOINT ["/app/inventory"]
