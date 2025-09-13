# gRPC-buf

A modern Golang service template featuring dual protocol APIs (gRPC/REST) using Connect, with deployment support for Google Cloud Run and event streaming capabilities.

[![Go Version](https://img.shields.io/github/go-mod/go-version/dipjyotimetia/gRPC-buf)](go.mod)
[![License](https://img.shields.io/github/license/dipjyotimetia/gRPC-buf)](LICENSE)

## Overview

gRPC-buf is a production-ready template for building cloud-native microservices with:

- 🚀 **Dual Protocol Support**: gRPC and REST endpoints via [connectrpc.com/connect](https://connectrpc.com)
- 🌩️ **Cloud-Native**: Ready for Google Cloud Run deployment
- 📨 **Event Streaming**: Built-in support for Pub/Sub and Kafka
- 🛠️ **Modern Tooling**: Streamlined development with Buf for Protocol Buffers
- 🔍 **Observability**: OpenTelemetry integration for monitoring
- 🗄️ **Persistence**: PostgreSQL database integration
- 💰 **Expenses API**: Resource-style CRUD for expenses

## Features

### API Development
- Unified gRPC and REST API endpoints using Connect
- Automatic OpenAPI documentation generation
- Protocol Buffer validation and linting with Buf
- Type-safe API contracts

### Cloud Integration
- Containerized deployment to Google Cloud Run
- Google Pub/Sub integration for event-driven architectures
- Kafka support for messaging
- Comprehensive OpenTelemetry instrumentation
- PostgreSQL for reliable data persistence

### Developer Experience
- Fast development workflow with hot reload
- Comprehensive test suite with integration tests
- Makefile-based task automation
- Docker containerization for consistent environments
- Database migrations

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

Expense API (REST examples):
- Create: POST `/v1/expenses` body: `{ "expense": { "user_id": "<uuid>", "amount": {"currency_code":"USD","units":"10"}, "category":"food" } }`
- Get: GET `/v1/expenses/{id}`
- List: GET `/v1/expenses?user_id=<uuid>`
- Update: PATCH `/v1/expenses/{id}` body: `{ "expense": {"id": "{id}", "description":"new"}, "update_mask": {"paths":["description"]}}`
- Delete: DELETE `/v1/expenses/{id}`

Environment variables:
- `ENVIRONMENT` (dev|prod) default: dev in `make run`
- `DATABASE_URL` Postgres connection (prod required)
- `DB_MAX_CONNS`/`DB_MIN_CONNS` tune pool sizes
- `OTEL_EXPORTER_OTLP_ENDPOINT` default: `otel-collector:4317`
- `OTEL_SERVICE_NAME` default: `grpc-buf`
- `JWT_SECRET` JWT signing key (prod required)

Configuration files:
- YAML configs live under `config/` and are loaded at startup:
  - `config/local.yaml` (used when `ENVIRONMENT` is dev or not set)
  - `config/production.yaml` (used when `ENVIRONMENT` is prod)
- Override with `CONFIG_PATH=/path/to/config.yaml`.
- Values in config are exported to environment variables for compatibility with existing code.

Advanced overrides with Koanf:
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

## Project Structure

```
.
├── cmd/                # Application entry points
│   ├── main.go         # Main application
│   └── server/         # Server implementation
├── internal/           # Private application code
│   ├── const/          # Constants
│   ├── gen/proto/      # Generated protocol buffer code
│   ├── logz/           # Logging utilities
│   ├── postgres/       # Database access layer
│   └── service/        # Business logic
├── proto/              # Protocol buffer definitions
│   ├── google/         # Google API definitions
│   ├── payment/        # Payment service definitions
│   ├── expense/        # Expense service definitions
│   └── registration/   # User registration definitions
└── scripts/            # Utility scripts
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
