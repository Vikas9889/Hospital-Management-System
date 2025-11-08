# Hospital Management System (HMS)

This repository contains a modular Hospital Management System implemented as 4 microservices in Go:

- **user-service** (patients & doctors) — port 8081
- **appointment-service** (bookings) — port 8082
- **billing-service** (bills & payments) — port 8083
- **notification-service** (mock notifications) — port 8084

Each service runs in Docker, each with its own PostgreSQL DB.

## Run
```bash
docker compose up --build
```
