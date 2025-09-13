# Architecture

This service exposes dual-protocol APIs (gRPC and REST) using Connect over a single HTTP/2 h2c port. It persists to PostgreSQL, ships traces to an OTel collector (Jaeger by default in docker-compose), and is containerized for Cloud Run.

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
    T[OTel Tracer]
  end

  A-->G
  B-->G
  G-->S
  S-->D
  S-->T

  subgraph Observability
    O[OTel Collector]
    J[Jaeger]
  end

  T-->O-->J
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
  participant OT as OTel

  C->>H: HTTP/2 Request
  H->>X: Route to handler
  X->>S: Typed request (Connect)
  activate S
  S->>OT: Start span
  S->>DB: Query/Exec
  DB-->>S: Rows/Result
  S->>OT: End span
  deactivate S
  S-->>X: Typed response
  X-->>H: Encoded (REST/gRPC)
  H-->>C: Response
```

Components
- Connect Handlers: `cmd/server/handler.go` wires service implementations to HTTP mux; also serves health and reflection.
- Service Layer: `internal/service` provides interfaces; implementations delegate to the datastore.
- Datastore (pgx pool): `internal/postgres` with embedded migrations and query methods.
- Configuration: `internal/config` (Koanf) loads YAML + env overrides and exports envs for compatibility.
- Observability: `internal/logz` creates an OTLP exporter; handlers add OTel middleware.

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
- OTel gRPC exporter: default `otel-collector:4317`

