# Architecture

## Request Flow

```
HTTP → Route → Middleware → Controller → Service → Repository → PostgreSQL
Response ← Serializer ← Service ← Repository ← DB
```

## Layers

| Layer        | Path                       | Job                                                                 |
| ------------ | -------------------------- | ------------------------------------------------------------------- |
| Route        | `modules/*/route.go`       | Endpoint registration, middleware attachment, dependency wiring    |
| Middleware   | `middleware/`              | CORS, JWTAuth, Recovery, RequireRoles                               |
| Controller   | `modules/*/controller.go`  | Bind/validate request, call service, return response               |
| Service      | `modules/*/service.go`     | Business rules, error → AppError mapping                            |
| Repository   | `modules/*/repository.go`  | SQL queries, row scanning, DB error detection                      |
| Model        | `modules/*/model/model.go` | DB entity structs (plain Go structs with `json:` tags)              |
| DTO          | `modules/*/dto/*.go`       | Request structs with `binding:` validator tags                      |
| Serializer   | `modules/*/serializer/*.go`| Response structs + model→response transformers                      |

## Module Structure (Standard)

```
src/modules/<feature>/
├── dto/          create_*, list_*, update_*_request.go
├── model/        model.go
├── serializer/   create_*, get_*, list_*, update_*_response.go
├── controller.go
├── repository.go
├── route.go
└── service.go
```

## Dependency Injection

Wired top-down at startup, no DI framework:
1. `main.go` → `database.Connect()` → `*sql.DB`
2. `router.New(db)` → passes `db` to each `Module.RegisterRoutes(r, db)`
3. `route.go`: `NewController(db)` → `NewRepository(db)` → `NewService(repo)`

Controller constructor pattern:
```go
func NewController(db *sql.DB) *Controller {
    return &Controller{service: NewService(NewRepository(db))}
}
```

## Interface-First

Service and Repository expose public interfaces; concrete structs (`service`, `postgresRepository`) are private lowercase. Constructors return the interface type (not `*service` directly), enabling testing/mocking.
