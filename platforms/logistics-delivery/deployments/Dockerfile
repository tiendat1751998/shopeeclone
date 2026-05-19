FROM golang:1.26-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/server ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/server /usr/local/bin/server
COPY --from=builder /app/configs /etc/logistics
COPY --from=builder /app/migrations /etc/logistics/migrations

EXPOSE 8080 9090
ENTRYPOINT ["server"]

