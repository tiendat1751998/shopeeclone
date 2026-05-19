# ---- Build stage ----
FROM golang:tip-trixie AS builder

RUN apt-get update && apt-get install -y git ca-certificates tzdata

WORKDIR /build

# Cache dependency downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o /build/server \
    ./cmd/server

# ---- Runtime stage ----
FROM scratch

# Copy CA certs and timezone data from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /build/server /server

# Copy migrations
COPY --from=builder /build/migrations /migrations

# Use non-root user (numeric for Kubernetes security contexts)
USER 65534:65534

# Expose ports
EXPOSE 8080 9090

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD ["/server", "-health-check"] || exit 1

ENTRYPOINT ["/server"]

