FROM golang:1.26.3-bookworm AS builder

RUN apt-get update && apt-get install -y git ca-certificates tzdata

WORKDIR /app/platforms/production-dashboard

COPY packages/go-shared /app/packages/go-shared

COPY platforms/production-dashboard/go.mod platforms/production-dashboard/go.sum ./
RUN go mod download
COPY platforms/production-dashboard/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o /app/bin/dashboard \
    ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app

COPY --from=builder /app/bin/dashboard /app/dashboard
COPY --from=builder /app/platforms/production-dashboard/migrations /app/migrations

EXPOSE 8080 9090

ENV GOMAXPROCS=0 \
    GOGC=50
ENTRYPOINT ["/app/dashboard"]
