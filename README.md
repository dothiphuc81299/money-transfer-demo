# 💸 Money Transfer Demo

A modular, microservices-based backend system for handling money transfers. Built in Go using the Gin framework and PostgreSQL, with a focus on maintainability, scalability, and clean architecture.

---

## 📐 System Overview

### 📁 Project Structure

```text
money-transfer-demo/
│
├── cmd/               # App entry points
│   ├── identity/      # Identity service (user, member management)
│   │   └── Dockerfile
│   ├── payment/       # Payment service (deposit, withdrawal)
│       └── Dockerfile
│
├── pkg/               # Shared and core logic
│   ├── identity/      # Identity domain logic
│   ├── payment/       # Payment domain logic (bank account, deposit, withdrawal, transfer)
│   ├── infra/         # Database (PostgreSQL) and other infrastructure
│   └── util/          # Common utilities
│
├── init-db.sql        # Database initialization script for Docker
├── docker-compose.yml # Docker Compose configuration
├── .github/workflows/ci-cd.yml # GitHub Actions CI/CD pipeline
└── README.md
```
---

### 🧠 System Architecture
![system_architecture](https://private-user-images.githubusercontent.com/48743415/442810730-d70cd844-4d79-41dc-ba00-bab84d8d4aad.jpg?jwt=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJnaXRodWIuY29tIiwiYXVkIjoicmF3LmdpdGh1YnVzZXJjb250ZW50LmNvbSIsImtleSI6ImtleTUiLCJleHAiOjE3NTg3OTIxMTIsIm5iZiI6MTc1ODc5MTgxMiwicGF0aCI6Ii80ODc0MzQxNS80NDI4MTA3MzAtZDcwY2Q4NDQtNGQ3OS00MWRjLWJhMDAtYmFiODRkOGQ0YWFkLmpwZz9YLUFtei1BbGdvcml0aG09QVdTNC1ITUFDLVNIQTI1NiZYLUFtei1DcmVkZW50aWFsPUFLSUFWQ09EWUxTQTUzUFFLNFpBJTJGMjAyNTA5MjUlMkZ1cy1lYXN0LTElMkZzMyUyRmF3czRfcmVxdWVzdCZYLUFtei1EYXRlPTIwMjUwOTI1VDA5MTY1MlomWC1BbXotRXhwaXJlcz0zMDAmWC1BbXotU2lnbmF0dXJlPTIwYTE2NTUxYThjMWJjZDQ3ODRkNWIyMDI2M2ZlYmIxOWFmNjQ1ZDY0MTYxZDNhN2NlYjFmNTliZTFiOThjNTkmWC1BbXotU2lnbmVkSGVhZGVycz1ob3N0In0.BfG0CYfB68jZc_gF5Lm_H7pJPrns6oJ8_AAXoHnRoqM)


---

## 📬 API Collection

You can try out the APIs using the Postman collection below:

👉 [Postman Collection](https://documenter.getpostman.com/view/12048946/2sB2jAbU18)

## 🚀 Getting Started

### 🔧 Prerequisites

- [Go 1.21+](https://go.dev/doc/install)
- [PostgreSQL](https://www.postgresql.org/download/)


### 🌀 Clone the repository

```bash
git clone https://github.com/dothiphuc81299/money-transfer-demo.git
cd money-transfer-demo
```

---
## 🛠️ Running Services
There are **two ways** to run the services: **Local Go** or **Docker Compose**.

---

### 1️⃣ Run with Docker Compose (recommended)

Docker Compose will automatically:

* Build images from `cmd/identity/Dockerfile` and `cmd/payment/Dockerfile`
* Create and initialize PostgreSQL databases using `init-db.sql`
* Start the services and map ports

```bash
docker-compose up --build
```

#### Ports

| Service          | Container Port | Host Port |
| ---------------- | -------------- | --------- |
| Identity Service | 5088           | 5088      |
| Payment Service  | 5090           | 5090      |
| PostgreSQL       | 5432           | 5433      |

> **Note**: `defaultPayment` in code should be `"payment-service:50004"` when using Docker Compose.

#### Health Checks

* Identity API: `http://localhost:5088/api/health`
* Payment API: `http://localhost:5090/api/health`

---

### 2️⃣ Run locally with Go

If you want to run without Docker:

1. Make sure PostgreSQL is running locally
2. **Create databases manually** (Docker init-db.sql is not used in this mode):

```sql
CREATE DATABASE identity-demo;
CREATE DATABASE payment-demo;
```

3. Run each service in separate terminals:

```bash
go run cmd/identity/main.go
go run cmd/payment/main.go
```

> **Note**: `defaultPayment` in config should be `"localhost:50004"` when running locally.

#### Ports

| Service          | Host Port |
| ---------------- | --------- |
| Identity Service | 5088      |
| Payment Service  | 5090      |

---

## 💳 Payment System

The Payment service supports three types of transactions:

1. **Local Bank Transfer**  
   Deposit and withdrawal of funds using a local bank transfer system.

2. **PayPal Integration**  
   Deposit and withdrawal of funds using PayPal.  
   Support for PayPal fee calculation and fee estimation.

3. **Account Transfer**  
   Transfer of funds between accounts within the system.

All three types of transactions (deposit, withdrawal, and transfer) are processed through the Payment service.

---




## 👤 Author

- [Phúc Đỗ](https://github.com/dothiphuc81299)

---

## 📝 License

MIT License. See [LICENSE](./LICENSE) for details

