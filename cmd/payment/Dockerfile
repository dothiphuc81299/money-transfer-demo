# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Copy migrations vào builder stage
COPY pkg/payment/db/migrations ./pkg/payment/db/migrations

# Build binary
RUN go build -o payment ./cmd/payment/main.go

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

# Copy binary từ builder
COPY --from=builder /app/payment .

# Copy migrations vào path mà Go code cần
COPY --from=builder /app/pkg/payment/db/migrations /pkg/payment/db/migrations

EXPOSE 5090

ENTRYPOINT ["./payment"]
