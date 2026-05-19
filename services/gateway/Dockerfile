FROM golang:tip-trixie AS builder

RUN apt-get update && apt-get install -y git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
COPY ../../packages/go-shared /app/packages/go-shared
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o /app/bin/gateway \
    ./cmd/server

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl

RUN addgroup -S gateway && adduser -S gateway -G gateway

WORKDIR /app

COPY --from=builder /app/bin/gateway /app/gateway
COPY --from=builder /app/configs /app/configs

USER gateway

EXPOSE 8080
EXPOSE 9090

HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/gateway"]

