FROM golang:tip-trixie AS builder
RUN apt-get update && apt-get install -y git ca-certificates tzdata
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/shipment-service ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /bin/shipment-service /app/shipment-service
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/configs /app/configs
RUN chown -R appuser:appgroup /app
USER appuser
EXPOSE 8085 9095
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 CMD wget --no-verbose --tries=1 --spider http://localhost:8085/health/live || exit 1
ENTRYPOINT ["/app/shipment-service"]

