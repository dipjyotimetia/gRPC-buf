# Operations

Running Locally
- Use docker-compose: `docker compose up -d`
  - Services: `api` (port 8080), `postgres` (5432), `otel-collector` (4317), `jaeger` (16686)
- Or run the app directly: `make run` (expects Postgres reachable per `config/local.yaml`)

Health & Readiness
- Liveness: `GET /livez` returns 200 when the server is up.
- gRPC Health: via grpchealth handler for service names.

Observability
- Tracing: OTel exporter to `OTEL_EXPORTER_OTLP_ENDPOINT` (default collector:4317)
- Jaeger UI: http://localhost:16686 (docker-compose)

Logging
- slog JSON to stdout.
- Level set by config (`server.log_level`) or env `LOG_LEVEL`.

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
- Cloud Run: see `.github/workflows/gcpdeploy.yaml`. Configure secrets and envs:
  - `DATABASE_URL`, `JWT_SECRET`, `OTEL_EXPORTER_OTLP_ENDPOINT`, `ENVIRONMENT=prod`.

Troubleshooting
- DB connection failures:
  - Ensure `DATABASE_URL` is valid and Postgres is reachable.
  - Check pool limits (`database.max_conns/min_conns`).
- Tracing not visible:
  - Verify collector endpoint and network.
  - Check service name `otel.service_name`.

