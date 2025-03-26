# Database Migration Guide

## 1️⃣ **Setup PostgreSQL Connection**
Ensure PostgreSQL is running and accessible. Update connection details in the `.env` file if needed.

## 2️⃣ **Run Migrations**
Run the following command to apply migrations:
```sh
migrate -database "postgres://postgres:secret@localhost:5432/identity-demo?sslmode=disable" -path pkg/identity/db/migrations up
```

## 3️⃣ **Check Migration Status**
Verify if the migration has been applied:
```sh
psql -h localhost -U postgres -d identity-demo -c "SELECT * FROM schema_migrations;"
```
If you see an entry with `version = 1`, the migration has been applied successfully.

## 4️⃣ **Verify Tables**
Check if the `member` table exists:
```sh
psql -h localhost -U postgres -d identity-demo -c "\dt public.*"
```
To inspect the `member` table:
```sh
psql -h localhost -U postgres -d identity-demo -c "SELECT * FROM public.member LIMIT 5;"
```

## 5️⃣ **Manually Apply Migration (if needed)**
If the table does not exist, run the SQL migration file manually:
```sh
psql -h localhost -U postgres -d identity-demo -f pkg/identity/db/migrations/000001_create_member_table.up.sql
```

## 6️⃣ **Reset and Reapply Migration (if needed)**
If migration fails or needs to be reapplied:
```sh
psql -h localhost -U postgres -d identity-demo -c "DELETE FROM schema_migrations WHERE version = 1;"
migrate -database "postgres://postgres:secret@localhost:5432/identity-demo?sslmode=disable" -path pkg/identity/db/migrations up
```

# Setting Up Protocol Buffers for Go

## Prerequisites
Ensure you have `protoc` installed. You can check with:
```sh
protoc --version
```
If not installed, download it from [Protocol Buffers releases](https://github.com/protocolbuffers/protobuf/releases) and follow the installation instructions.

## Install Go Plugins
Run the following commands to install the necessary plugins:
```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Verify Installation
Check if the plugins are installed correctly:
```sh
which protoc-gen-go
which protoc-gen-go-grpc
```
Ensure they are located in `$GOPATH/bin` (default: `~/go/bin`).

## Add to PATH (if necessary)
If `which protoc-gen-go` returns nothing, add the Go binary directory to your `PATH`:
```sh
export PATH=$PATH:$(go env GOPATH)/bin
```
Add this line to `~/.bashrc`, `~/.bash_profile`, or `~/.zshrc` to persist it:
```sh
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

## Generating Go Code from .proto Files
Run the following command to generate Go code:
```sh
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       pkg/apis/payment/payment.proto
```

Replace `pkg/apis/payment/payment.proto` with your `.proto` file path as needed.

## Troubleshooting
- If `protoc-gen-go: program not found or is not executable` appears, ensure you installed `protoc-gen-go` and added it to `PATH`.
- If `Unknown flag: --go-grpc_opt` appears, check that `protoc-gen-go-grpc` is correctly installed.

Now you're ready to work with Protocol Buffers in Go! 🚀

