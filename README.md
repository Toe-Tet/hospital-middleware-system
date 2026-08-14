# Hospital Middleware System

REST API middleware for multi-hospital deployments, built with Go, Gin, and PostgreSQL.

## Overview

This project follows a layered module structure:

- Route
- Middleware
- Controller
- Service
- Repository
- PostgreSQL

Current domain modules:

- `hospitals`
- `staffs`
- `patients`

The application also includes JWT authentication, Swagger documentation, SQL migrations, and Go-based seeders.

## Tech Stack

- Go `1.26.5`
- Gin `v1.12.0`
- PostgreSQL via `database/sql` + `pgx/v5`
- JWT via `golang-jwt/jwt/v5`
- Validation via `go-playground/validator/v10`
- Logging via `zerolog`
- Swagger via `swaggo/gin-swagger`
- Migrations via `golang-migrate`

## Project Structure

```text
.
├── cmd/api/                   # application entrypoint
├── docs/                      # generated Swagger artifacts
├── scripts/                   # migration / rollback / seed helpers
├── src/
│   ├── config/                # environment config loading
│   ├── database/              # postgres connection, migrations, seeders
│   ├── errors/                # AppError and status mapping
│   ├── helper/                # JWT, validation, response helpers
│   ├── logger/                # zerolog setup
│   ├── middleware/            # CORS, recovery, JWT auth
│   ├── modules/               # hospitals, staffs, patients
│   └── router/                # top-level route registration
├── .env.example
├── Makefile
└── go.mod
```

## Prerequisites

- Go `1.26.5`
- PostgreSQL running locally or reachable from this service
- A database created for the app, for example `hospital_middleware`

## Configuration

Copy the example file and adjust values for your environment:

```bash
cp .env.example .env
```

Environment variables:

| Variable            | Default                                          | Description             |
| ------------------- | ------------------------------------------------ | ----------------------- |
| `APP_ENV`           | `development`                                    | Application environment |
| `APP_PORT`          | `8080`                                           | HTTP server port        |
| `DB_HOST`           | `localhost`                                      | PostgreSQL host         |
| `DB_PORT`           | `5432`                                           | PostgreSQL port         |
| `DB_USER`           | `postgres`                                       | PostgreSQL user         |
| `DB_PASSWORD`       | `postgres`                                       | PostgreSQL password     |
| `DB_NAME`           | `hospital_middleware`                            | Database name           |
| `DB_SSLMODE`        | `disable`                                        | PostgreSQL SSL mode     |
| `JWT_SECRET`        | `change-me-to-a-long-random-string-min-32-chars` | JWT signing secret      |
| `JWT_EXPIRES_HOURS` | `24`                                             | Token lifetime in hours |

## Getting Started

1. Install local developer tools:

```bash
make install-tools
```

2. Create and configure `.env`.

3. Apply migrations:

```bash
make migrate-up
```

4. Seed sample data:

```bash
make seed
```

5. Start the API:

```bash
make run
```

The service starts on `http://localhost:8080` by default.

## Common Commands

| Command                 | Purpose                                     |
| ----------------------- | ------------------------------------------- |
| `make build`            | Build the binary into `./bin/`              |
| `make run`              | Run the API locally                         |
| `make test`             | Run all tests                               |
| `make test-coverage`    | Generate coverage report                    |
| `make fmt`              | Run `gofmt` and `goimports`                 |
| `make tidy`             | Run `go mod tidy`                           |
| `make install-tools`    | Install local `migrate` and `swag` binaries |
| `make migrate-up`       | Apply all pending migrations                |
| `make migrate-down`     | Roll back the last migration                |
| `make migrate-down-all` | Roll back all migrations                    |
| `make seed`             | Run Go seeders                              |
| `make swagger`          | Regenerate Swagger docs                     |

## API Docs

- Health check: [http://localhost:8080/health](http://localhost:8080/health)
- Swagger UI: [http://localhost:8080/api/docs/index.html](http://localhost:8080/api/docs/index.html)

## Current Routes

Routes below reflect the current code registration in `src/router/` and each module `route.go`.

### Public Routes

| Method | Path             | Description          |
| ------ | ---------------- | -------------------- |
| `GET`  | `/health`        | Service health check |
| `GET`  | `/api/docs/*any` | Swagger docs         |
| `GET`  | `/hospitals`     | List hospitals       |
| `POST` | `/staffs`        | Create staff         |
| `POST` | `/staffs/create` | Create staff alias   |
| `POST` | `/staffs/login`  | Staff login          |

### JWT-Protected Routes

| Method | Path            | Description           |
| ------ | --------------- | --------------------- |
| `GET`  | `/staffs/me`    | Current staff profile |
| `GET`  | `/patients`     | List patients         |
| `GET`  | `/patients/:id` | Get patient by ID     |

## Seed Data

After `make seed`, sample hospitals, staffs, and patients are inserted.

Example staff credentials from the current seeder:

- Username: `admin`
- Password: `password`

Other seeded usernames include:

- `dr-smith`
- `carol`
- `dave`
- `dr-johnson`
- `nurse-frank`

## Data Models

### Hospital

```text
{ ID, Name, Address, Phone, Status, CreatedAt, UpdatedAt }
```

### Staff

```text
{ ID, HospitalID, Username, Name(*string), Password(json:"-"), Status, CreatedAt(*time.Time), UpdatedAt(*time.Time) }
```

### Patient

```text
{ ID, HospitalID, PatientHN, FirstNameTh(*string), MiddleNameTh(*string), LastNameTh(*string), FirstNameEn, MiddleNameEn(*string), LastNameEn, NationalID(*string), PassportID(*string), DateOfBirth, Gender, Email, PhoneNumber, Status, CreatedAt, UpdatedAt(*time.Time) }
```

## Architecture Notes

- Entry point: `cmd/api/main.go`
- Router setup: `src/router/router.go`
- Modules live under `src/modules/<feature>/`
- Repositories use raw SQL with PostgreSQL
- Responses use helper wrappers such as `helper.OK()` and `helper.Error()`
- JWT claims include `staff_id`, `hospital_id`, and `role`

## Testing

Run the full test suite:

```bash
make test
```

Or run directly with Go:

```bash
go test ./...
```

## Known Implementation Notes

- The README route list reflects the current registered routes, even where internal specs may describe broader or stricter behavior.
- `patients` endpoints are JWT-protected.
- `staffs/me` is JWT-protected.
- `hospitals` list is currently exposed without JWT middleware.
- Internal project notes indicate some module specs are ahead of implementation, especially around RBAC and tenancy checks.
