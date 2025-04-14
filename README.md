# Database Migration Guide


## 2️⃣ **Run Migrations**
Run the following command to apply migrations:
```sh
migrate -database "postgres://postgres:secret@localhost:5432/identity-demo?sslmode=disable" -path pkg/identity/db/migrations up
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



## Generating Go Code from .proto Files
Run the following command to generate Go code:
```sh
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       pkg/apis/payment/payment.proto
```
