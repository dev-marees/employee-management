# Employee Management System — Backend API

A production-grade REST API for an Employee Management System, built with **Go 1.25**, **Gin**, **GORM**, and **PostgreSQL**, following **Clean Architecture** and **SOLID** principles.

It powers a React frontend with: employee listing, search, filter, sort, CRUD, soft delete, dashboard statistics, and JWT-based authentication with role-based authorization.

---

## Features

- Clean Architecture: `handler → service → repository → model` with DTO boundaries
- JWT authentication (access + refresh tokens, independently signed)
- Role-based authorization: `Admin`, `HR`, `Employee`
- Employee CRUD with **soft delete**
- Pagination, search, filtering, and whitelisted sorting
- Dashboard aggregation (totals + department breakdown)
- Request validation with field-level error messages
- Structured logging (zap) with request IDs
- Graceful shutdown
- Swagger / OpenAPI documentation
- Dockerfile + docker-compose
- Unit tests for the JWT, pagination, and service layers

---

## Architecture

```
cmd/api/main.go              # entrypoint: config, logger, db, server, graceful shutdown
internal/
  app/                       # composition root (dependency injection)
  router/                    # route registration + middleware wiring
  config/                    # env-based configuration (viper)
  database/                  # GORM connection + migration helpers
  logger/                    # zap logger factory
  middleware/                # auth, RBAC, request-id, logging, recovery, CORS
  auth/        {handler, service, repository, model, dto}
  employee/    {handler, service, repository, model, dto}
  dashboard/   {handler, service, repository}
pkg/
  apperror/                  # transport-agnostic sentinel errors
  jwt/                       # token issuing/verification
  hash/                      # bcrypt password hashing
  pagination/                # pagination params + result envelope
  response/                  # success/error JSON envelopes
  validation/                # validator -> field error mapping
docs/                        # generated Swagger spec
migrations/                  # versioned SQL migrations
```

**Dependency rule:** dependencies point inward. Handlers depend on service
interfaces; services depend on repository interfaces; nothing in the domain
imports Gin or GORM-specific transport concerns. This keeps each layer unit-
testable with in-memory fakes (see `*_test.go`).

---

## Getting Started

### Prerequisites
- Go 1.25+
- PostgreSQL 14+ (or Docker)

### 1. Configure environment
```bash
cp .env.example .env
# edit JWT secrets and DB credentials
```

### 2. Run with Docker Compose (recommended)
```bash
make up          # builds the API image and starts Postgres + API
# API:     http://localhost:8080
# Swagger: http://localhost:8080/swagger/index.html
```

### 3. Run locally (Go + your own Postgres)
```bash
make tidy        # resolve dependencies
make run         # starts on :8080
```

The service auto-migrates the schema on boot via GORM. For production, prefer
the versioned SQL files in `migrations/` (run with `golang-migrate`).

---

## API

Base path: `/api/v1`

| Method | Endpoint              | Auth        | Roles            |
|--------|-----------------------|-------------|------------------|
| POST   | `/auth/register`      | public      | —                |
| POST   | `/auth/login`         | public      | —                |
| POST   | `/auth/refresh`       | public      | —                |
| GET    | `/dashboard`          | Bearer      | any              |
| GET    | `/employees`          | Bearer      | any              |
| GET    | `/employees/:id`      | Bearer      | any              |
| POST   | `/employees`          | Bearer      | Admin, HR        |
| PUT    | `/employees/:id`      | Bearer      | Admin, HR        |
| DELETE | `/employees/:id`      | Bearer      | Admin            |

### Query parameters for `GET /employees`
```
?page=1&limit=10
&search=john            # matches code, first/last name, email, department
&department=Engineering
&status=active          # active | inactive
&joined_from=2023-01-01&joined_to=2023-12-31
&sort_by=salary         # name | salary | joining_date
&sort_direction=desc    # asc | desc
```

### Response envelopes

Success:
```json
{ "success": true, "message": "Employee created successfully", "data": {} }
```

Paginated list (`data` field):
```json
{ "data": [], "page": 1, "limit": 10, "total": 100, "total_pages": 10 }
```

Error:
```json
{ "success": false, "message": "Validation failed", "errors": { "email": "Must be a valid email address" } }
```

### Example: register, login, then create an employee
```bash
# Register an HR user
curl -X POST localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"HR Admin","email":"hr@acme.com","password":"supersecret","role":"HR"}'

# Login -> copy access_token from data
TOKEN=$(curl -s -X POST localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"hr@acme.com","password":"supersecret"}' | jq -r .data.access_token)

# Create an employee
curl -X POST localhost:8080/api/v1/employees \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"employee_code":"EMP-001","first_name":"John","last_name":"Smith","email":"john@acme.com","department":"Engineering","salary":95000,"joining_date":"2023-04-01"}'
```

---

## Development

```bash
make help          # list all targets
make test          # run unit tests
make test-cover    # tests + coverage summary
make lint          # go vet
make fmt           # gofmt -s
make swagger       # regenerate docs (needs `go install github.com/swaggo/swag/cmd/swag@latest`)
```

---

## Notes & Production Hardening

- **CORS** is permissive (`*`) by default — restrict `Access-Control-Allow-Origin`
  in `internal/middleware/common.go` for production.
- **Refresh tokens** are stateless here. For revocation, persist a token/jti
  allowlist or rotate with a stored hash.
- Run **versioned migrations** (`migrations/`) in production instead of relying
  on `AutoMigrate`.
- Set strong, distinct `JWT_ACCESS_SECRET` / `JWT_REFRESH_SECRET` values.
