FROM golang:1.26.3-bookworm AS builder

RUN apt-get update && apt-get install -y git ca-certificates tzdata

WORKDIR /app/services/shipment

COPY packages/go-shared /app/packages/go-shared
COPY services/shipment/go.mod services/shipment/go.sum ./
RUN go mod download

COPY services/shipment/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o /app/bin/shipment \
    ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app

COPY --from=builder /app/bin/shipment /app/shipment
COPY --from=builder /app/services/shipment/migrations /app/migrations

EXPOSE 8080 9090

ENTRYPOINT ["/app/shipment"]
