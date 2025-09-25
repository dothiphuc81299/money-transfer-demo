# ========================
# Stage 1: Build
# ========================
FROM golang:1.23-alpine AS builder

# Thư mục làm việc trong container builder
WORKDIR /app

# COPY file quản lý dependency trước để tận dụng cache Docker
COPY go.mod .
COPY go.sum .
RUN go mod download

# COPY toàn bộ source code + migrations từ host → builder stage
# Stage builder → compile Go binary, tận dụng cache, image lớn
COPY . .
COPY pkg/identity/db/migrations ./pkg/identity/db/migrations

# Build binary Go từ source code
RUN go build -o identity ./cmd/identity/main.go


# ========================
# Stage 2: Run
# ========================
FROM alpine:latest

# Thư mục làm việc trong container runtime
WORKDIR /app

# COPY từ builder stage → runtime stage
# Stage runtime → chỉ chạy binary + migrations → image nhỏ, không cần Go toolchain
COPY --from=builder /app/identity .
COPY --from=builder /app/pkg/identity/db/migrations /pkg/identity/db/migrations
# Nếu bỏ copy migrations → container sẽ lỗi khi chạy migration

# Metadata port, không bắt buộc, giúp Docker/Docker-compose biết port container lắng nghe
EXPOSE 5088

# ENTRYPOINT bắt buộc, container start sẽ tự động chạy binary
ENTRYPOINT ["./identity"]

# Notes:
# 1. Stage builder: compile Go binary, tận dụng cache, nhẹ nhàng
# 2. Stage runtime: chỉ chạy binary + migrations, image nhỏ, không cần Go toolchain
# 3. COPY migrations → bắt buộc nếu ứng dụng cần migration
# 4. EXPOSE → metadata, không ảnh hưởng việc chạy container
# 5. ENTRYPOINT → container start tự động chạy binary
#docker build -f cmd/identity/Dockerfile -t identity .