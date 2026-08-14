# Conventions

## Errors (`src/errors/`)

Use ONLY `AppError` constructors in services — never raw errors to controllers:

- `NewBadRequest(msg)` → 400
- `NewValidationError(msg, details)` → 400
- `NewUnauthorized(msg)` → 401
- `NewForbidden(msg)` → 403
- `NewNotFound(msg)` → 404
- `NewConflict(msg)` → 409
- `NewInternal(err)` → 500 (wrap raw DB/driver errors here)

**Repository layer** returns raw errors or short sentinel strings; service maps to AppError.

## Responses (`src/helper/response.go`)

Envelope shape:

```json
{ "success": true|false, "data": {}, "error": { "code": "...", "message": "...", "details": {} } }
```

Controller helpers ONLY — no direct `c.JSON()`:

- `helper.OK(c, data)` 200
- `helper.NoContent(c)` 204
- `helper.Paginated(c, items, total, page, perPage)` 200 (items/total/page/per_page/total_pages)
- `helper.Error(c, err)` → auto status from AppError

## Validation (Controller)

Two-step pattern, always:

1. `c.ShouldBindJSON(&req)`
2. `helper.ValidateStruct(&req)`

DTO fields use `binding:"required,min=N,max=N,email"` tags (go-playground/validator v10).

## Repository (SQL)

- Placeholders: `$1, $2, ...` (never sprintf/concat into SQL)
- Context variants: `QueryRowContext`, `QueryContext`, `ExecContext`
- NotFound (GetByID): `errors.Is(err, sql.ErrNoRows)` → return `nil, nil` (service → 404)
- Unique violation: detect `"23505"` (PG code) → return sentinel string (service → `NewConflict`)
- Lists: separate `COUNT(*)` query + `LIMIT $1 OFFSET $2`; `defer rows.Close()` right after `QueryContext`; check `rows.Err()` after loop
- Timestamps: use `NOW()` in SQL, not Go time

## Auth & RBAC

**JWT** (`middleware.JWTAuth()`):

- Requires `Authorization: Bearer <token>`
- Validates token; injects into ctx: `staff_id`, `hospital_id`, `role`
- Access via: `middleware.GetStaffID(c)`, `middleware.GetHospitalID(c)`
- Token signed with `JWT_SECRET`, expiry from `JWT_EXPIRES_HOURS` (default 24h)
- Claims: `StaffID`, `HospitalID`, `Role`

**Passwords**: bcrypt (`bcrypt.DefaultCost`) — hash on Create, compare on Login. `Password` field has `json:"-"`.

**Route groups** (`router.go`):

- `protected := v1.Group(""); protected.Use(JWTAuth())` — default home for all modules
- Public endpoints (login, health) outside this group

**Role gating**: per-route `middleware.RequireRoles("admin", ...)` before controller handler.

## Routing

Health: `GET /health` (public).

| Action | Pattern                 | Status        |
| ------ | ----------------------- | ------------- |
| Create | `POST /resources`       | 201           |
| List   | `GET /resources`        | 200 paginated |
| Get    | `GET /resources/:id`    | 200           |
| Update | `PUT /resources/:id`    | 200           |
| Delete | `DELETE /resources/:id` | 204           |

Parse ID params via `helper.ParseIDParam(c, "id")` → validates int, returns BadRequest on failure.

## Pagination

Query params: `page` (default 1), `per_page` (default 10). Parsed via `helper.ParsePageParams(c)`.
Offset: `(page - 1) * perPage`. Repo returns `([]*Model, int64 total, error)`.

## Logging & Config

**Logging** (zerolog):

- Init: `logger.Init(config.AppConfig.AppEnv)`; global `logger.Log`
- `gin.Logger()` also attached for requests
- Dev = pretty mode; prod = JSON
- DB connect failure is the only Fatal; runtime returns AppError

**Config**:

- Loaded first via `config.Load()` — uses `.env` only if `APP_ENV` unset
- Access via singleton `config.AppConfig`
- Keys: `APP_ENV`, `APP_PORT`, `DB_{HOST,PORT,USER,PASSWORD,NAME,SSLMODE}`, `JWT_SECRET`, `JWT_EXPIRES_HOURS`
- `cfg.DBConnectionString()` (libpq key=value for `sql.Open`); `cfg.DBURL()` (postgres:// URL for migrate CLI)
