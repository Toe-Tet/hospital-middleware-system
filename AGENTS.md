# Hospital Middleware System - Agent Guidelines

## Role: Senior Backend Engineer (Go Expert)

You are a senior Go backend engineer with deep expertise in Go idioms, Gin, PostgreSQL (raw `database/sql` + pgx/v5, no ORM), REST API design, JWT auth + RBAC, and maintainable layered architecture.

## Project Summary

Hospital Middleware System — REST API middleware for multi-hospital deployments managing Hospitals, Staff (JWT login), and Patients. Designed for extensibility (appointments, etc.) via a consistent module pattern.

## Tech Stack

| Concern    | Library / Tool                          | Version  |
| ---------- | --------------------------------------- | -------- |
| Language   | Go                                      | 1.26.5   |
| Web        | `gin-gonic/gin`                         | v1.12.0  |
| Database   | `jackc/pgx/v5` (via `database/sql`)     | v5.10.0  |
| Validation | `go-playground/validator/v10`           | v10.30.3 |
| JWT        | `golang-jwt/jwt/v5`                     | v5.3.1   |
| Password   | `golang.org/x/crypto/bcrypt`            | v0.55.0  |
| Config     | `joho/godotenv`                         | v1.5.1   |
| Logging    | `rs/zerolog`                            | v1.35.1  |
| Migrations | `golang-migrate/migrate/v4` (local CLI) | v4.19.1  |

## Specification Index

Detailed specs live under `.trae/specs/`. Read the relevant spec before coding.

| File                                                         | What's inside                                                                                                                        |
| ------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------ |
| [`.trae/specs/architecture.md`](.trae/specs/architecture.md) | Request flow, layer responsibilities, standard module directory structure, DI wiring, interface-first pattern                        |
| [`.trae/specs/conventions.md`](.trae/specs/conventions.md)   | Errors (AppError constructors), response envelope, validation, repo SQL rules, JWT auth + RBAC, routing, pagination, logging, config |
| [`.trae/specs/modules.md`](.trae/specs/modules.md)           | Hospitals, Staffs, Patients — endpoints, models, constraints, known implementation gaps, FK relationships                            |
| [`.trae/specs/database.md`](.trae/specs/database.md)         | Connection + pool config, migration rules + commands, seeders, PG error codes (23505 etc.), query safety                             |
| [`.trae/specs/development.md`](.trae/specs/development.md)   | Makefile targets, code style rules, file naming, security checklist, feature addition checklist, graceful shutdown                   |

## Quick Principles

- **Strict layers**: Controller → Service (interface) → Repository (interface) → PostgreSQL. No short-circuits.
- **Interface-first**: Public Service / Repository interfaces; private `service` / `postgresRepository` structs.
- **AppError in services only**: Use `NewBadRequest/NewValidationError/NewUnauthorized/NewForbidden/NewNotFound/NewConflict/NewInternal`. Repos return raw/sentinel errors.
- **Raw SQL only, no ORM**: `$1, $2` placeholders; `QueryRowContext`/`QueryContext`; `defer rows.Close()`; detect PG `23505` → conflict sentinel.
- **JWT claims**: `staff_id`, `hospital_id`.
- **Response helpers only**: `helper.OK/NoContent/Paginated/Error` — no direct `c.JSON()`.
- **Validation two-step**: `c.ShouldBindJSON(&req)` then `helper.ValidateStruct(&req)`; DTOs use `binding:` tags.
- **No comments in Go source** (project convention).
- **Before feature done**: `make fmt tidy` → `make build` → `make test` → migrations + smoke tests.

## Entry Points

| Concern    | Path                                                              |
| ---------- | ----------------------------------------------------------------- |
| App main   | `cmd/api/main.go`                                                 |
| Router     | `src/router/router.go`                                            |
| Config     | `src/config/config.go` (env via `.env`)                           |
| DB         | `src/database/postgres.go`                                        |
| Migrations | `src/database/migrations/`                                        |
| Seeders    | `src/database/seeders/` (build tag `seeder`)                      |
| Middleware | `src/middleware/` (CORS, Recovery, JWTAuth)                       |
| Helpers    | `src/helper/` (response, JWT, validation, pagination, ID parsing) |
| Errors     | `src/errors/` (AppError + constructors + status mapper)           |
| Modules    | `src/modules/{hospitals,staffs,patients}/`                        |
