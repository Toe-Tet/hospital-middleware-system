# Development

## Makefile Targets

| Category | Target               | Action                                                              |
| -------- | -------------------- | ------------------------------------------------------------------- |
| Build    | `make build`         | Compile → `./bin/hospital-middleware`                               |
| Run      | `make run`           | `go run ./cmd/api/main.go`                                          |
| Clean    | `make clean`         | Remove `./bin/`, `coverage.out`, `coverage.html`                    |
| Test     | `make test`          | `go test -v ./...`                                                  |
| Test     | `make test-coverage` | Coverage profile + HTML report                                       |
| DB       | `make install-tools` | Install `golang-migrate` → `./bin/`                                 |
| DB       | `make migrate-up`    | Apply pending migrations                                            |
| DB       | `make migrate-down`  | Rollback last migration                                             |
| DB       | `make migrate-down-all` | Rollback ALL (dev only — destructive)                           |
| DB       | `make seed`          | Run Go seeders (tag: seeder)                                        |
| Quality  | `make fmt`           | `gofmt -s` + `goimports` (best effort)                              |
| Quality  | `make tidy`          | `go mod tidy`                                                       |

## Code Style

1. **No comments in Go code** (project convention).
2. Import groups (3 blocks, blank line separated): stdlib → 3rd-party → internal `hospital-middleware-system/...`
3. Wrap errors with `%w`: `fmt.Errorf("open db: %w", err)`
4. Pass `ctx context.Context` as first param to all I/O functions; propagate to all DB calls
5. Services/repos: public Interface, private lowercase struct
6. `Password` always `json:"-"`; never log/serialize JWTs or PII
7. New entities default `Status = "active"`; timestamps use SQL `NOW()`
8. Empty query results → return `[]*Entity{}` not `nil` (JSON `[]` vs `null`)
9. Nullable DB columns → pointer type in model (`*time.Time`, `*string`)

## Naming Files

| Kind                     | Pattern                             | Example                               |
| ------------------------ | ----------------------------------- | ------------------------------------- |
| DTO request              | `<action>_<entity>_request.go`      | `create_hospital_request.go`          |
| Serializer response      | `<action>_<entity>_response.go`     | `get_hospital_response.go`            |
| Migration                | `<NNNNNN>_<snake_name>.[up\|down].sql` | `000004_create_appointments.up.sql` |
| Seeder                   | `<entity>_seeder.go`                | `hospital_seeder.go`                  |

## Security Checklist

- [ ] Parameterized SQL everywhere (no sprintf into query)
- [ ] Protected endpoints inside `JWTAuth()` group; public endpoints (login, `/health`) outside
- [ ] Write/admin routes wrapped with `RequireRoles("admin"[,...])`
- [ ] Passwords bcrypt only; never plaintext logged or stored
- [ ] DTOs carry validator tags; controllers run 2-step validation
- [ ] Service layer enforces hospital tenancy (`hospital_id` from JWT claim matches queried resource) unless system-admin role
- [ ] Login failure always 401 + generic "Invalid email or password" (never leak "email not found")

## Feature Addition Checklist

1. Add migration pair (up + down) with columns, constraints, indexes
2. Add model struct (pointers for NULLable)
3. Add DTO structs with validator tags (create/list/update)
4. Add serializer structs + `Serialize*()` funcs
5. Implement repository interface + `postgresRepository`; handle ErrNoRows→nil, 23505→sentinel
6. Implement service interface + AppError mapping + business rules/tenancy
7. Controller with bind→validate→service→response helpers (no direct `c.JSON`)
8. Wire `route.go`; apply `RequireRoles` where needed
9. Register module in `router.go` (protected/public group)
10. Optional: add seeder + hook into `seeder.go` run list after dependencies
11. `make fmt tidy` → `make build` → `make test`
12. Run migrations + seed + smoke-test endpoints locally

## Graceful Shutdown

Server already handles SIGINT/SIGTERM: 15s shutdown timeout; in-flight requests drain before DB `Close()`.

If you add background goroutines (cron, consumers), tie their lifecycle to a cancellable context derived in `main.run`; ensure they exit within the 15s shutdown window or risk ungraceful termination.
