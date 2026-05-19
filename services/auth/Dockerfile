FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
COPY ../../packages/go-shared /app/packages/go-shared
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o /app/bin/auth \
    ./cmd/server

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl

RUN addgroup -S auth && adduser -S auth -G auth

WORKDIR /app

COPY --from=builder /app/bin/auth /app/auth
COPY --from=builder /app/migrations /app/migrations

USER auth

EXPOSE 8080
EXPOSE 9090

HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/auth"]

