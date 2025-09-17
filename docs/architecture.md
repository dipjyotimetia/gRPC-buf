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

Data Model (partial)

```mermaid
erDiagram
  USERS {
    uuid id PK
    text email
    text password
    text first_name
    text last_name
    timestamptz created_at
    timestamptz updated_at
  }

  PAYMENTS {
    uuid id PK
    bigint card_no
    int card_type
    text name
    text address
    real amount
    timestamptz created_at
  }

  EXPENSES {
    uuid id PK
    uuid user_id
    bigint amount_cents
    text currency_code
    text category
    text description
    timestamptz created_at
    timestamptz updated_at
  }
```

Ports & Protocols
- App: HTTP/2 h2c on port 8080 (gRPC + REST via Connect)
