# Build stage
FROM golang:tip-trixie AS builder

RUN apt-get update && apt-get install -y git ca-certificates tzdata

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o /build/bin/product \
    ./cmd/server

# Runtime stage
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app

COPY --from=builder /build/bin/product /app/product
COPY --from=builder /build/migrations /app/migrations
COPY --from=builder /build/configs /app/configs

EXPOSE 8080 9090

ENTRYPOINT ["/app/product"]