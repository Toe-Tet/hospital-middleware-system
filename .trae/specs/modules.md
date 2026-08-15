# Modules

All modules mounted via `router.go`. Default group is `protected` (requires JWT).

## Hospitals (`/hospitals`)

| Method | Path         | Auth   | Action           |
| ------ | ------------ | ------ | ---------------- |
| GET    | `/hospitals` | Public | List (paginated) |

**Model**: `{ ID, Name, Address, Phone, Status, CreatedAt, UpdatedAt }`

- Status default: `"active"` (CHECK `active`/`inactive` enforced by PG)
- `name` is UNIQUE → 409 Conflict via PG `23505` detection

## Staffs (`/staffs`)

Wired endpoints:
| Method | Path | Auth | Action |
| ------ | -------- | ---- | ------ |
| POST | `/staffs` | JWT | Create staff |

Controller methods exist but NOT yet wired in route.go / router.go:
| Method | Path | Auth | Action |
| ------ | --------------- | ------- | ------------------------ |
| POST | `/staffs/login` | Public | Credentials → JWT token |
| GET | `/staffs/me` | JWT | Self info from claims |
| GET | `/staffs/:id` | JWT | Fetch staff by ID |

**Model**: `{ ID, HospitalID, Username, Name(*string nullable), Password(json:"-"), Status, CreatedAt(*time.Time), UpdatedAt(*time.Time nullable) }`

- Status default: `"active"`; login checks `Status == "active"` → else 401 "Account is inactive"
- `username` UNIQUE → 409 Conflict
- `hospital_id` FK validated on create → 400 "hospital does not exist"
- Password via bcrypt
- Login response: `{ Token, ExpiresAt, Staff(Me shape) }`
- JWT claims: `StaffID`, `HospitalID`

Known gaps to fix when wiring:

- `Login` must be outside `JWTAuth()` protected group

## Patients (`/patients`)

| Method | Path            | Auth | Action           |
| ------ | --------------- | ---- | ---------------- |
| GET    | `/patients`     | JWT  | List (paginated) |
| GET    | `/patients/:id` | JWT  | Get by ID        |

**Model**: `{ ID, HospitalID, PatientHN, FirstNameTh(*string nullable), MiddleNameTh(*string nullable), LastNameTh(*string nullable), FirstNameEn, MiddleNameEn(*string nullable), LastNameEn, NationalID(*string nullable), PassportID(*string nullable), DateOfBirth, Gender, Email, PhoneNumber, Status, CreatedAt, UpdatedAt(*time.Time nullable) }`

Known gap: **List/Get do not scope by requester's `hospital_id`**.

## Relationships

```
hospitals (1) ──< staffs (N via hospital_id FK)
hospitals (1) ──< patients (N via hospital_id FK)
```
