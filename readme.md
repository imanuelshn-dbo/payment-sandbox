Payment Simulation Sandbox API

Tech Stack :

- Golang (Gin)
- PostgreSQL (Main DB)
- MongoDB (Logging)
- GORM (ORM)
- JWT Authentication
- Swagger (API Docs)


Project Structure :

payment-sandbox/
├─cmd/app               # main entrypoint
├── internal/
│   ├── config           # DB & config
│   ├── middleware       # JWT, Role, Idempotency
│   ├── modules          # auth, wallet, invoice, payment, refund, admin
│   ├── models           # DB models
│   ├── logger           # MongoDB logging
├── pkg/
│   ├── response         # standard API response
│   ├── app-error        # custom error



Setup Environment :

Buat file .env di root project:
APP_PORT=8080

DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=your_user
DB_PASS=your_password
DB_NAME=payment_service
DB_SSLMODE=disable

JWT_SECRET=supersecretkey

MONGO_URI=mongodb://localhost:27017
MONGO_DB=payment_logs


Setup Database :

PostgreSQL :
Install PostgreSQL
CREATE DATABASE payment_service;

MongoDB :
Install MongoDB atau gunakan MongoDB Compass
Pastikan berjalan di:
mongodb://localhost:27017

Run Application :

Masuk ke folder
cd cmd/app

Jalankan:
go run main.go

Server akan berjalan di:
http://localhost:8080
