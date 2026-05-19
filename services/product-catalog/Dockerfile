FROM golang:tip-trixie AS builder
WORKDIR /app; COPY go.mod go.sum ./; RUN go mod download; COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/server
FROM gcr.io/distroless/static-debian12; COPY --from=builder /app/main /main; USER nonroot:nonroot; EXPOSE 8080 9090; ENTRYPOINT ["/main"]

