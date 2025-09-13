# Operations

Running Locally
- Use docker-compose: `docker compose up -d`
  - Services: `api` (port from `config/local.yaml`), `postgres` (5432)
- Or run the app directly: `make run` (expects Postgres reachable per `config/local.yaml`)

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
- Cloud Run: see `.github/workflows/gcpdeploy.yaml`. Provide configuration via `CONFIG_PATH` or bake config into image. Provide secrets (e.g., `DATABASE_URL`, `JWT_SECRET`) via Secret Manager/env.

Troubleshooting
- DB connection failures:
  - Ensure `DATABASE_URL` is valid and Postgres is reachable.
  - Check pool limits (`database.max_conns/min_conns`).
 
