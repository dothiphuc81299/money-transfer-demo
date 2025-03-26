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

