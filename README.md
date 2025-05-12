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
│   ├── payment/       # Payment service (deposit, withdrawal,transfer)
│
├── pkg/               # Shared and core logic
│   ├── identity/      # Identity domain logic
│   ├── payment/       # Payment domain logic (bank account, deposit, withdrawal, transfer)
│   ├── infra/         # Database (PostgreSQL) and other infrastructure
│   └── util/          # Common utilities
```

### 🧠 System Architecture

![system-architecture](https://github.com/user-attachments/assets/d70cd844-4d79-41dc-ba00-bab84d8d4aad)

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

### 🛠️ Create the Database

Before running the services, ensure the `payment-demo` and `identity-demo` database is created in your PostgreSQL instance. You can do this by running the following SQL command in your PostgreSQL client or terminal:

```sql
CREATE DATABASE payment-demo;
CREATE DATABASE identity-demo;
```

## 🛠️ Run the services

In separate terminals, run each service manually:

1. **Identity Service** (manages users and members)
```bash
go run cmd/identity/main.go
```
2. **Payment Service** (manages deposit, withdrawal, transfer)

```bash
go run cmd/payment/main.go
```

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

MIT License. See [LICENSE](./LICENSE) for detail

