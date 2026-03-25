# Meet Sushruta - Hospital Management System Backend

A clean architecture Go backend for a hospital management system built with:
- **Go 1.22**
- **Gin** - HTTP web framework
- **GORM** - ORM for database operations
- **PostgreSQL** - Database
- **UUID** - Primary keys for all models

## Project Structure

```
meet_sushruta_backend/
├── config/              # Configuration management
│   ├── config.go       # Environment config with Viper
│   └── database.go     # GORM initialization & migrations
├── model/              # Data models (15 GORM structs)
│   ├── user.go
│   ├── patient.go
│   ├── doctor.go
│   ├── nurse.go
│   ├── doctor_schedule.go
│   ├── appointment.go
│   ├── vitals.go
│   ├── medicine.go
│   ├── prescription.go
│   ├── prescription_item.go
│   ├── lab_request.go
│   ├── bill.go
│   ├── bed.go
│   ├── audit_log.go
│   └── notification.go
├── repository/         # Data access layer (to be implemented)
├── service/            # Business logic layer (to be implemented)
├── handler/            # HTTP request handlers (to be implemented)
├── middleware/         # HTTP middleware (to be implemented)
├── main.go            # Application entry point
├── go.mod             # Go module definition
├── go.sum             # Go dependencies checksums
├── .env.example       # Example environment variables
└── .gitignore
```

## Database Models

All 15 models include:
- UUID primary keys via `github.com/google/uuid`
- Embedded `gorm.Model` (CreatedAt, UpdatedAt, DeletedAt)
- **Exception**: `AuditLog` - immutable, no DeletedAt field
- Proper GORM tags with foreign keys and constraints
- Indexes on frequently queried fields

### Models

1. **User** - Base user with role-based access (admin, patient, doctor, nurse)
2. **Patient** - Patient profile linked to User
3. **Doctor** - Doctor profile linked to User with specialization
4. **Nurse** - Nurse profile linked to User with shift management
5. **DoctorSchedule** - Doctor availability schedule (day/time slots)
6. **Appointment** - Links Patient ↔ Doctor with status tracking (pending, confirmed, in_progress, completed, cancelled)
7. **Vitals** - Patient vital signs (temperature, BP, heart rate, etc.)
8. **Medicine** - Medicine catalog with stock management
9. **Prescription** - Prescriptions issued by Doctor to Patient
10. **PrescriptionItem** - Individual medicine items in a prescription
11. **LabRequest** - Lab test requests with status tracking
12. **Bill** - Patient billing with payment status
13. **Bed** - Hospital bed management and occupancy
14. **AuditLog** - Immutable audit trail (no DeletedAt)
15. **Notification** - User notifications with delivery tracking

## Configuration

Environment variables (see `.env.example`):

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=meet_sushruta
DB_USER=postgres
DB_PASSWORD=your_secure_password
DB_SSLMODE=disable

# Server
SERVER_PORT=8080

# JWT
JWT_SECRET=your_jwt_secret_key_change_in_production
JWT_EXPIRY_MINUTES=60
REFRESH_EXPIRY_DAYS=7
```

## Setup Instructions

### 1. Install PostgreSQL
Ensure PostgreSQL is installed and running.

### 2. Create Database
```sql
CREATE DATABASE meet_sushruta;
```

### 3. Clone & Setup Repository
```bash
cd meet_sushruta_backend
go mod tidy
```

### 4. Configure Environment
```bash
cp .env.example .env
# Edit .env with your database credentials
```

### 5. Run Application
```bash
go run main.go
```

The application will:
- Load configuration from environment variables
- Connect to PostgreSQL
- Run all GORM AutoMigrate for 15 models
- Start Gin server on port 8080 (or configured port)

## Architecture

This project follows **Clean Architecture** principles:

- **Model Layer** - Database entities with relationships
- **Repository Layer** - Database access abstraction (to be implemented)
- **Service Layer** - Business logic and rules (to be implemented)
- **Handler Layer** - HTTP request/response handling (to be implemented)
- **Middleware Layer** - Cross-cutting concerns like auth, logging (to be implemented)

## Next Steps

1. Implement repository interfaces for data access
2. Implement service layer with business logic
3. Create handlers for HTTP endpoints
4. Add authentication middleware (JWT)
5. Add request validation
6. Add error handling middleware
7. Create API routes
8. Add unit and integration tests

## Dependencies

- `github.com/gin-gonic/gin` - Web framework
- `github.com/google/uuid` - UUID generation
- `github.com/spf13/viper` - Configuration management
- `gorm.io/gorm` - ORM
- `gorm.io/driver/postgres` - PostgreSQL driver
