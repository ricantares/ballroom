# Ballroom API Architecture

## Overview

Ballroom is a RESTful API built with Go using the Gin web framework. It provides a school management system with role-based access control (RBAC) and JWT authentication.

## Project Structure

```
ballroom/
├── cmd/web/              # Application entry point
├── internal/
│   ├── api/              # HTTP handlers and routing
│   ├── db/               # Data access layer (Repository pattern)
│   ├── domain/           # Business entities and models
│   ├── logger/           # Logging utilities
│   ├── rest/             # HATEOAS response formatting
│   └── security/         # JWT & RBAC authentication
├── mocks/                # Mock implementations for testing
├── resources/            # Static assets and configuration
└── api_test/            # API contract tests
```

---

## Layer Architecture

### 1. Entry Point — `cmd/web/main.go`

The application starts here, initializing:

- **Environment variables** — loads `.env` file
- **Timezone configuration** — sets application timezone
- **Logger** — initializes file-based logging
- **RBAC** — loads role-based access control rules from `rbac.json`
- **Database** — establishes PostgreSQL connection
- **HTTP Server** — configures and starts the Gin router

### 2. API Layer — `internal/api/`

| File | Responsibility |
|------|----------------|
| `route.go` | Defines HTTP endpoints and applies middleware |
| `handler.go` | Contains request/response handling logic |
| `login.go` | Authentication endpoint |
| `scuola.go` | School and room (sala) management endpoints |
| `wrapper.go` | Middleware wrappers (e.g., RBAC enforcement) |
| `home.go` | Home page handler |

**Key Patterns:**
- Route handlers implement interfaces defined in `handler.go`
- RBAC middleware protects authenticated routes
- Responses use HATEOAS format for hypermedia linking

### 3. Domain Layer — `internal/domain/`

Core business entities:

| Entity | Description |
|--------|-------------|
| `model.go` | Base `Model` struct with `Id`, `Created_at`, `Updated_at`, `Deleted` |
| `scuola.go` | School entity |
| `utente.go` | User entity with roles |
| `persona.go` | Person (generic contact info) |
| `corso.go` | Course entity |
| `gara.go` | Competition entity |
| `insegnante.go` | Teacher entity |
| `affiliazione.go` | Affiliation entity |
| `certificato.go` | Certificate entity |
| `contatto.go` | Contact details |

**Key Type:** `Uuid` — custom type alias for `int32` representing unique identifiers

### 4. Data Access Layer — `internal/db/`

Implements the **Repository Pattern**:

```
┌─────────────────┐
│   API Handler   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Repository     │  (internal/db/repo.go)
│  - GetScuola()  │
│  - GetSala()    │
│  - ListSala()   │
│  - CreateSala() │
│  - ...          │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   DB Interface │  (internal/db/repo.go)
│   - pgx driver │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  PostgreSQL     │
└─────────────────┘
```

**Key Files:**
- `repo.go` — Repository struct and methods
- `postgres.go` — PostgreSQL connection and query execution
- `scuolasvc.go` — School-specific data operations
- `utentesvc.go` — User-specific data operations

### 5. Security Layer — `internal/security/`

#### JWT Authentication (`jwt.go`)

- Generates JWT tokens with custom claims including user role
- Token includes: `ExpiresAt`, `IssuedAt`, `Issuer`, `Subject`, `Role`
- Secret key loaded from environment variables

#### RBAC (`rbac.go`)

Role-based access control loaded from `rbac.json`:

```json
{
  "rbac": [
    {
      "pattern": "/scuola/sala/\\d*",
      "permissions": {
        "can_read": ["Admin", "Direzione", "Iscritto", "Maestro", "Staff"],
        "can_write": ["Admin"],
        "can_delete": ["Admin"]
      }
    }
  ]
}
```

- **Pattern** — Regex matching requested resource
- **Permissions** — Maps HTTP methods to allowed roles:
  - `can_read` → GET
  - `can_write` → POST, PUT
  - `can_delete` → DELETE

### 6. REST Layer — `internal/rest/`

HATEOAS (Hypermedia as the Engine of Application State) implementation:

| Component | Purpose |
|-----------|---------|
| `hateoas.go` | Builds HAL-compliant responses with links and embedded resources |
| `Link` struct | `{ Rel, Href }` pairs for hypermedia |
| `BaseResource` | Container for links and data |
| `Page` struct | Pagination metadata (items, size, number, pages) |

**Response Format:**
```json
{
  "_links": {
    "self": { "href": "..." }
  },
  "data": { ... }
}
```

---

## Request Flow

```
┌──────────────┐
│   Client     │
└──────┬───────┘
       │ HTTP Request
       ▼
┌──────────────────────────────────────┐
│  Gin Router                          │
│  - route.go defines endpoints        │
└──────┬───────────────────────────────┘
       │
       ▼ (if protected)
┌──────────────────────────────────────┐
│  RBAC Middleware                     │
│  - wrapper.go: RbacHandler()         │
│  - Validates JWT + checks role       │
└──────┬───────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────┐
│  Handler                              │
│  - login.go, scuola.go, etc.         │
│  - Calls Repository methods         │
└──────┬───────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────┐
│  Repository                           │
│  - repo.go: wraps DB interface        │
└──────┬───────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────┐
│  Database (PostgreSQL)                │
│  - pgx driver                        │
└──────────────────────────────────────┘
       │
       ▼ (response)
┌──────────────────────────────────────┐
│  HATEOAS Formatter                    │
│  - rest/hateoas.go adds _links       │
└──────┬───────────────────────────────┘
       │
       ▼
┌──────────────┐
│   Client     │
└──────────────┘
```

---

## Technology Stack

| Component | Technology |
|-----------|-------------|
| Web Framework | Gin (github.com/gin-gonic/gin) |
| Database Driver | pgx (github.com/jackc/pgx/v5) |
| JWT | golang-jwt/jwt (github.com/golang-jwt/jwt/v5) |
| HATEOAS | go2hal (github.com/pmoule/go2hal) |
| Configuration | godotenv (github.com/joho/godotenv) |
| Testing | testify (github.com/stretchr/testify) |
| Go Version | 1.23.1 |

---

## Configuration Files

| File | Purpose |
|------|---------|
| `.env` | Environment variables (database, JWT secret, timezone) |
| `rbac.json` | Role-based access control rules |
| `go.mod` | Go module dependencies |

---

## Key Design Decisions

1. **Single School Model** — The application manages a single school entity, created at initialization and immutable
2. **Repository Pattern** — Data access is abstracted through the Repository interface, enabling testability and flexibility
3. **HATEOAS Responses** — All API responses include hypermedia links for discoverability
4. **RBAC via Middleware** — Authorization is enforced at the router level before reaching handlers
5. **Custom UUID Type** — Domain uses `domain.Uuid` (int32) instead of standard UUID for simplicity