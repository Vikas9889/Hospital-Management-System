# Hospital Management System

A modern, microservices-based hospital management system built with Go, Docker, and Kubernetes.

## System Architecture

The system consists of the following microservices:

1. **User Service** (Port: 8081)
   - Manages patient and staff profiles
   - Handles user authentication and role management
   - PII data protection with masking for logs

2. **Doctor Service** (Port: 8085)
   - Manages doctor profiles and specializations
   - Handles doctor availability and schedules
   - Tracks doctor experience and departments

3. **Appointment Service** (Port: 8082)
   - Manages patient appointments
   - Handles scheduling and rescheduling
   - Tracks appointment status (SCHEDULED, COMPLETED, CANCELLED)

4. **Billing Service** (Port: 8083)
   - Generates bills for appointments
   - Manages payment status
   - Handles currency and amount calculations

5. **Notification Service** (Port: 8084)
   - Sends notifications for appointments
   - Handles reminders and updates
   - Supports multiple notification channels

6. **Prescription Service** (Port: 8087)
   - Manages medical prescriptions
   - Records medications and dosages
   - Links prescriptions with appointments

7. **Payment Service** (Port: 8088)
   - Processes payments
   - Supports multiple payment methods
   - Handles payment status updates

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.19 or later
- PostgreSQL client (optional, for direct DB access)
- PowerShell 5.1 or later (for running scripts)
- Kubernetes/Minikube (for k8s deployment)

### Local Development Setup

1. **Clone the Repository**
   ```powershell
   git clone https://github.com/Vikas9889/Hospital-Management-System.git
   cd Hospital-Management-System
   ```

2. **Start Services**
   ```powershell
   docker compose up -d --build
   ```

3. **Run Demo Flow**
   ```powershell
   ./scripts/demo-flow.ps1
   ```

### API Endpoints

#### User Service (8081)
- POST `/v1/users` - Create new user
- GET `/v1/users/{id}` - Get user details
- GET `/v1/users` - List users
- PUT `/v1/users/{id}` - Update user

#### Doctor Service (8085)
- POST `/v1/doctors` - Create doctor profile
- GET `/v1/doctors/{id}` - Get doctor details
- GET `/v1/doctors` - List doctors
- PUT `/v1/doctors/{id}` - Update doctor profile

#### Appointment Service (8082)
- POST `/v1/appointments` - Create appointment
- GET `/v1/appointments/{id}` - Get appointment details
- POST `/v1/appointments/{id}/complete` - Mark appointment as completed
- GET `/v1/appointments` - List appointments

#### Billing Service (8083)
- POST `/v1/bills` - Generate new bill
- GET `/v1/bills/{id}` - Get bill details
- GET `/v1/bills` - List bills
- PUT `/v1/bills/{id}` - Update bill status

#### Prescription Service (8087)
- POST `/v1/prescriptions` - Create prescription
- GET `/v1/prescriptions` - List prescriptions
- GET `/v1/prescriptions/{id}` - Get prescription details

#### Payment Service (8088)
- POST `/v1/payments` - Process payment
- GET `/v1/payments` - List payments
- GET `/v1/payments/{id}` - Get payment details

### Database Schema

Each service maintains its own PostgreSQL database:

- users_db (5432)
- doctors_db (5435)
- appointments_db (5433)
- billing_db (5434)

### Kubernetes Deployment

1. **Build Images**
   ```powershell
   docker compose build
   ```

2. **Deploy to Kubernetes**
   ```powershell
   kubectl apply -f k8s/all.yaml
   ```

3. **Verify Deployment**
   ```powershell
   kubectl get pods,svc -A
   ```

### Monitoring & Health Checks

- Each service exposes a `/health` endpoint
- Use `docker compose logs [service]` for service logs
- Kubernetes health probes configured in k8s/all.yaml

### Security Features

1. **PII Protection**
   - Email and phone masking in logs
   - Encrypted data transmission
   - Role-based access control

2. **API Security**
   - Input validation
   - Error handling
   - Rate limiting

### Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

### License

This project is licensed under the MIT License - see the LICENSE file for details.

### Contact

For questions and support, please open an issue in the GitHub repository. (HMS)

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
