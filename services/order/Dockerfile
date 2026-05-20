# Build stage
FROM golang:tip-trixie AS builder

RUN apt-get update && apt-get install -y git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o /app/bin/order \
    ./cmd/server

# Runtime stage
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app

COPY --from=builder /app/bin/order /app/order
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/configs /app/configs

EXPOSE 8083 9093

ENTRYPOINT ["/app/order"]