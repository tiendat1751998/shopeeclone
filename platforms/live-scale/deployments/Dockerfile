FROM golang:1.26.3-bookworm AS builder

RUN apt-get update && apt-get install -y gcc libc-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/server ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/server /usr/local/bin/server
COPY --from=builder /app/configs /etc/live-scale
COPY --from=builder /app/migrations /etc/live-scale/migrations

EXPOSE 8081 9091
ENTRYPOINT ["server"]

