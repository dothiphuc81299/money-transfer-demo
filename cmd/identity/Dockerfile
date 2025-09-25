# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
COPY pkg/identity/db/migrations ./pkg/identity/db/migrations

RUN go build -o identity ./cmd/identity/main.go

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/identity .
COPY --from=builder /app/pkg/identity/db/migrations /pkg/identity/db/migrations

EXPOSE 5088

ENTRYPOINT ["./identity"]
