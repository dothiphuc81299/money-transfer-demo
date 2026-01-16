# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o payment ./cmd/payment/main.go

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/payment .
COPY --from=builder /app/pkg/payment/db/migrations /pkg/payment/db/migrations

EXPOSE 5090

ENTRYPOINT ["./payment"]
