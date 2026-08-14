# Database

## Connection

- Driver: `jackc/pgx/v5/stdlib` via standard `database/sql` (no ORM)
- `database.Connect()` → `sql.Open("pgx", connStr)` then `PingContext` with 10s timeout
- Pool: MaxOpen=25, MaxIdle=10, MaxLifetime=5min, MaxIdleTime=5min
- Global: `database.DB` (\*sql.DB singleton); `database.Close()` deferred from main
- Config: `DBConnectionString()` (key=value for sql.Open), `DBURL()` (postgres:// for migrate)

## Migrations (golang-migrate v4)

**Setup**: `make install-tools` — downloads `./bin/migrate` (local, project-scoped)

**Files**: `src/database/migrations/`

- Naming: `<snake_case_name>.up.sql` + `.down.sql` (zero-padded numbering, eg. `000004_create_appointments.up.sql`)
- **Always paired** up/down; down must fully reverse changes
- Use explicit columns + indexes; never `SELECT *` in application queries

**Schema rules**:

- PK: `SERIAL` (or `BIGSERIAL`)
- Timestamps: `TIMESTAMPTZ NOT NULL DEFAULT NOW()` for `created_at`, `updated_at`
- Status enums: `VARCHAR(N) NOT NULL DEFAULT 'active' CHECK(status IN (...))`
- Index FK columns explicitly (PG doesn't auto); UNIQUE via named `CREATE UNIQUE INDEX`; add indexes on filtered/ordered columns

**Commands**:
| Command | Effect |
| ----------------------- | --------------------------------------- |
| `make migrate-up` | Run all pending |
| `make migrate-down` | Rollback last |
| `make migrate-down-all` | Rollback all (destructive, dev only) |

## Seeders

- Go-based, build tag `seeder`; entry `src/database/seeders/seeder.go`
- Run via `make seed` → `go run -tags=seeder ./src/database/seeders/...`
- Order (FKs): Hospitals → Staffs → Patients
- Default admin after seed: `admin@hospital.com` / `password` (dev only)

## Error Codes (Repository Layer)

| PG Code | Name             | Detect with `strings.Contains(err.Error(), "XXXX")` → Map to |
| ------- | ---------------- | ------------------------------------------------------------ |
| 23505   | unique_violation | sentinel string → service `NewConflict(msg)`                 |
| 23503   | FK violation     | `NewBadRequest` / `NewNotFound` depending on context         |
| 23514   | CHECK violation  | 400 if user-fixable; else 500                                |

Helper pattern:

```go
func isUniqueViolation(err error) bool {
    return err != nil && strings.Contains(err.Error(), "23505")
}
```

## Query Safety

- Always parameterized (`$1, $2…`) — never string-build SQL
- Use `*Context` variants; propagate request `ctx` through all layers
- `defer rows.Close()` immediately after successful `QueryContext`; check `rows.Err()` after loop
- Multi-statement writes: use `db.BeginTx(ctx, nil)`, defer `tx.Rollback()`, explicit `tx.Commit()`
- Single-row inserts: use `INSERT ... RETURNING ...` + `QueryRowContext.Scan` (one round-trip)
