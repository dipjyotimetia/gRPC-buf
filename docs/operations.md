# Operations

Running Locally
- Use docker-compose: `docker compose up -d`
  - Services: `api` (port from `config/local.yaml`), `postgres` (5432)
- Or run the app directly: `make run` (expects Postgres reachable per `config/local.yaml`)

## Testing

- Run unit tests: `make test`
- Run integration tests: `make test-integration` (requires a running server)

## Code Generation

- To regenerate code from the protobuf definitions, run: `make generate`

## MCP Server

This project includes an "MCP" server that communicates over standard I/O. This server is likely intended for use in a specific context, such as being run as a child process by another application. For most local development, you should use `make run` to start the main API server.

- To run the MCP server: `make mcp-run`

Health & Readiness
- Liveness: `GET /livez` returns 200 when the server process is running.
- Readiness: `GET /readyz` returns 200 when DB is reachable; 503 otherwise.
- Version: `GET /version` returns JSON with version, commit, and date.
- gRPC Health: via grpchealth handler for service names.

Observability
- Tracing: removed — no OpenTelemetry instrumentation is included.

Logging
- slog JSON to stdout.
- Level set by config (`server.log_level`).

Migrations
- Embedded migrations run on startup when `server.run_migrations = true`.
- Manual:
  - Up: `make migrate-up` (requires `DATABASE_URL`)
  - Down: `make migrate-down`

Security
- JWT signing key from `security.jwt_secret`. Required in production.
- CORS: set `server.cors_allowed_origins` (use exact origins in prod).

Deploy
- Docker: `docker build -t <image> .`
- Cloud Run: see `.github/workflows/gcpdeploy.yaml`. Provide configuration via `CONFIG_PATH` or bake config into image. Provide secrets (e.g., `DATABASE_URL`, `JWT_SECRET`) via Secret Manager/env. The service reads `PORT` from the environment.

Production
- The API binary fails fast in `ENVIRONMENT=prod` if configuration cannot be loaded or validated. Ensure either:
  - a valid `CONFIG_PATH` is provided, or
  - all required env vars are set (e.g., `DATABASE_URL`, `SECURITY_JWT_SECRET`).

Troubleshooting
- DB connection failures:
  - Ensure `DATABASE_URL` is valid and Postgres is reachable.
  - Check pool limits (`database.max_conns/min_conns`).
  - If deploying to Cloud Run and you see hostname resolving errors like `lookup postgres: no such host`,
    verify your `DATABASE_URL` doesn't reference the docker-compose hostname `postgres`.
    That hostname only works inside `docker compose` networks. Use a reachable host such as a Cloud SQL
    private/public IP or DNS name. Example: `postgres://USER:PASS@127.0.0.1:5432/DB?sslmode=require` (replace host appropriately).
 
