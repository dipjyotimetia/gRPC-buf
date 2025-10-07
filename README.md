# gRPC-buf

A modern Golang service template featuring dual protocol APIs (gRPC/REST) using Connect, with PostgreSQL persistence and production-focused defaults.

[![Go Version](https://img.shields.io/github/go-mod/go-version/dipjyotimetia/gRPC-buf)](go.mod)
[![License](https://img.shields.io/github/license/dipjyotimetia/gRPC-buf)](LICENSE)

## Overview

gRPC-buf is a lean, production-friendly service with:

- Dual protocol APIs via Connect (REST + gRPC on one port)
- PostgreSQL integration with embedded migrations
- JWT auth verification and simple login/registration flows
- Basic rate limiting for the login endpoint
- Health, readiness, and version endpoints
- Protobufs managed with Buf (lint, generate)

## Features

### What's Included
- Connect handlers for REST and gRPC
- MCP (Model Context Protocol) server for AI/LLM integration
- Buf-based generation and linting
- Integration tests (Docker compose-based)
- Makefile tasks (build, test, lint, migrations)
- Docker/Compose for local development

## Prerequisites

- Go 1.25 or later
- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- [Buf CLI](https://docs.buf.build/installation)
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install) (for deployment)

## Quick Start

1. **Clone the repository**:
   ```bash
   git clone https://github.com/dipjyotimetia/gRPC-buf.git
   cd gRPC-buf
   ```

2. **Setup local environment**:
   ```bash
   make setup
   ```

3. **Generate code from Protocol Buffers**:
   ```bash
   make generate
   ```

4. **Start the service**:
   ```bash
   make run
   ```

5. **Access the API**:
   - REST API: http://localhost:8080
   - gRPC: localhost:8080
    - Health: http://localhost:8080/livez
    - Readiness: http://localhost:8080/readyz
    - Version: http://localhost:8080/version

Expense API (REST examples):
- Create: POST `/v1/expenses` body: `{ "expense": { "user_id": "<uuid>", "amount": {"currency_code":"USD","units":"10"}, "category":"food" } }`
- Get: GET `/v1/expenses/{id}`
- List: GET `/v1/expenses?user_id=<uuid>`
- Update: PATCH `/v1/expenses/{id}` body: `{ "expense": {"id": "{id}", "description":"new"}, "update_mask": {"paths":["description"]}}`
- Delete: DELETE `/v1/expenses/{id}`

Environment variables (optional):
- `CONFIG_PATH` to point to a YAML config file.
- Production configs may reference environment variables (YAML interpolation) for secrets like `DATABASE_URL` and `JWT_SECRET`.

Configuration files:
- YAML configs live under `config/` and are loaded at startup:
  - `config/local.yaml` (default for local dev)
  - `config/production.yaml` (example for prod; secrets via environment interpolation if desired)
- Override file path with `CONFIG_PATH=/path/to/config.yaml`.

Advanced overrides (envconfig):
- Override any YAML key via environment variables. Preferred names mirror the struct hierarchy: `SERVER_PORT`, `DATABASE_URL`, etc.
- Legacy names with the `CFG_` prefix still work (e.g., `CFG_SERVER_PORT`).
- Examples:
  - `SERVER_PORT=9090` overrides `server.port`.
  - `CFG_DATABASE_MAX_CONNS=200` overrides `database.max_conns`.
  - `CFG_SERVER_LOG_LEVEL=debug` remains supported for backward compatibility.

Notes:
- Database pool sizing: Prefer `DATABASE_MAX_CONNS`/`DATABASE_MIN_CONNS`. Legacy `DB_MAX_CONNS`/`DB_MIN_CONNS` are still recognized.
- CORS: In dev, an empty list allows all. In any env, adding `"*"` to `server.cors_allowed_origins` allows all.

## Common Development Commands

| Command | Description |
|---------|-------------|
| `make test` | Run all tests |
| `make lint` | Run linters for Go code and Protocol Buffers |
| `make build` | Build the service binary |
| `make clean` | Clean up generated files |
| `make migrate-up` | Apply database migrations |
| `make migrate-down` | Roll back database migrations |
| `make migrate-run` | Run embedded migrations via Go |
| `make migrate-run-local DSN=postgres://...` | Run embedded migrations against a custom DSN |
| `make mcp-build` | Build MCP server binary |
| `make mcp-run` | Run MCP server locally |
| `make mcp-generate` | Generate MCP stubs from protobuf |

## MCP Server

The project includes an MCP (Model Context Protocol) server that exposes all gRPC services as MCP tools for AI/LLM integration.

### Running the MCP Server

```bash
# Build the MCP server
make mcp-build

# Run the MCP server
make mcp-run
```

### Available MCP Tools

The MCP server exposes the following services as tools:

- **Expense Service**: `CreateExpense`, `GetExpense`, `ListExpenses`, `UpdateExpense`, `DeleteExpense`
- **User Service**: `RegisterUser`, `LoginUser`
- **Payment Service**: `MakePayment`, `MarkInvoicePaid`, `PayInvoice`

### Using with MCP Clients

The MCP server uses stdio transport by default, making it compatible with MCP-enabled clients:

```json
{
  "mcpServers": {
    "grpc-buf": {
      "command": "/path/to/bin/mcp-server",
      "env": {
        "DATABASE_URL": "postgres://postgres:postgres@localhost:5432/grpcbuf?sslmode=disable",
        "CONFIG_PATH": "./config/local.yaml"
      }
    }
  }
}
```

## Project Structure

```
.
├── cmd/                    # Application entrypoints
│   ├── api/                # REST/gRPC service binary
│   │   └── main.go         # Main
│   └── mcp-server/         # MCP server binary
│       └── main.go         # MCP entrypoint
├── internal/               # Private application code
│   ├── config/             # Config loading & env export (envconfig)
│   ├── gen/proto/          # Generated protocol buffer code
│   │   ├── expense/        # Expense protos + MCP stubs
│   │   ├── payment/        # Payment protos + MCP stubs
│   │   └── registration/   # User protos + MCP stubs
│   ├── postgres/           # Database access layer + migrations
│   ├── security/           # JWT verification
│   ├── server/             # Server lifecycle (listen/shutdown, CORS, h2c)
│   ├── service/            # Service layer
│   │   └── mcp/            # MCP adapters for services
│   └── transport/          # Transport layers
│       ├── http/           # HTTP wiring, health, reflection
│       └── mcp/            # MCP server implementation
├── proto/              # Protocol buffer definitions
│   ├── google/         # Google API definitions
│   ├── payment/        # Payment service definitions
│   ├── expense/        # Expense service definitions
│   └── registration/   # User registration definitions
└── scripts/                # Utility scripts
```

## Documentation

- Docs index: docs/README.md
- Architecture: docs/architecture.md
- Workflows (CI/CD + Dev): docs/workflows.md
- Configuration: docs/configuration.md
- APIs: docs/apis.md
- Operations: docs/operations.md

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.
Migrations (local)
- Using embedded migrator with custom connection:
  - `make migrate-run-local DSN=postgres://user:pass@localhost:5432/grpcbuf?sslmode=disable`
- Using external migrator (requires migrate CLI):
  - `make migrate-up DSN=postgres://...`
  - `make migrate-down DSN=postgres://...`
Using .env for local dev
- Copy `.env.example` to `.env` and adjust values.
- The Makefile auto-loads `.env` when present.
- Example:
  - `CONFIG_PATH=./config/local.yaml`
  - `DATABASE_URL=postgres://postgres:postgres@localhost:5432/grpcbuf?sslmode=disable`
  - `JWT_SECRET=change-me`
