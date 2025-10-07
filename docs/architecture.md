# Architecture

This service exposes dual-protocol APIs (gRPC and REST) using Connect over a single HTTP/2 h2c port. It persists to PostgreSQL and ships with docker/compose for local and CI.

High-level

```mermaid
flowchart LR
  subgraph Clients
    A[Browser / Mobile]-- REST -->G
    B[gRPC Clients]-- gRPC -->G
  end

  subgraph Service
    G[HTTP/2 h2c Server\nConnect Handlers]
    S[Service Layer]
    D[(Postgres)]
  end

  A-->G
  B-->G
  G-->S
  S-->D
  

 
```

Request Flow

```mermaid
sequenceDiagram
  autonumber
  participant C as Client (REST/gRPC)
  participant H as HTTP Server (h2c)
  participant X as Connect Handler
  participant S as Service
  participant DB as Postgres

  C->>H: HTTP/2 Request
  H->>X: Route to handler
  X->>S: Typed request (Connect)
  activate S
  S->>DB: Query/Exec
  DB-->>S: Rows/Result
  deactivate S
  S-->>X: Typed response
  X-->>H: Encoded (REST/gRPC)
  H-->>C: Response
```

Components
- Connect Handlers: `internal/transport/http/handler.go` wires service implementations to HTTP mux; also serves health and reflection.
- Service Layer: `internal/service` provides interfaces; implementations delegate to the datastore.
- Datastore (pgx pool): `internal/postgres` with embedded migrations and query methods.
- Configuration: `internal/config` (envconfig) loads YAML + env overrides.
- Auth & Rate Limit: `internal/security` and `internal/transport/middleware/*`.

Services
- **UserService**: Manages user registration and login.
- **PaymentService**: Handles payments and invoices.
- **ExpenseService**: Manages expenses.

Data Model

```mermaid
erDiagram
  USERS {
    string id PK
    string email
    string password
    string first_name
    string last_name
    timestamp created_at
    timestamp updated_at
  }

  EXPENSES {
    string id PK
    string user_id FK
    Money amount
    string category
    string description
    timestamp create_time
    timestamp update_time
  }

  PAYMENTS {
    int64 card_no PK
    CardType card
    string name
    string address_lines
    float amount
    timestamp payment_created
  }

  USERS ||--o{ EXPENSES : has
```

Ports & Protocols
- App: HTTP/2 h2c on port 8080 (gRPC + REST via Connect)
