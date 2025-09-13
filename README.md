# gRPC-buf

A modern Golang service template featuring dual protocol APIs (gRPC/REST) using Connect, with PostgreSQL persistence and production-focused defaults.

[![Go Version](https://img.shields.io/github/go-mod/go-version/dipjyotimetia/gRPC-buf)](go.mod)
[![License](https://img.shields.io/github/license/dipjyotimetia/gRPC-buf)](LICENSE)

## Overview

gRPC-buf is a lean, production-friendly service starter with:

- Dual protocol APIs via Connect (REST + gRPC on one port)
- PostgreSQL integration with embedded migrations
- JWT auth verification and simple login/registration flows
- Basic rate limiting for the login endpoint
- Health, readiness, and version endpoints
- Protobufs managed with Buf (lint, generate)

## Features

### What’s Included
- Connect handlers for REST and gRPC
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

Advanced overrides (Koanf):
- You can override any YAML key via environment variables using the `CFG_` prefix and `__` as a nesting separator.
- Examples:
  - `CFG_SERVER__PORT=9090` overrides `server.port`.
  - `CFG_DATABASE__MAX_CONNS=200` overrides `database.max_conns`.
  - `CFG_SERVER__LOG_LEVEL=debug` sets the log level.

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

## Project Structure

```
.
├── cmd/                    # Application entrypoints
│   └── api/                # Service binary
│       └── main.go         # Main
├── internal/               # Private application code
│   ├── config/             # Config loading & env export (Koanf)
│   ├── gen/proto/          # Generated protocol buffer code
│   ├── postgres/           # Database access layer + migrations
│   ├── security/           # JWT verification
│   ├── server/             # Server lifecycle (listen/shutdown, CORS, h2c)
│   ├── service/            # Thin service layer delegating to datastore
│   └── transport/          # HTTP wiring, health, reflection
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
