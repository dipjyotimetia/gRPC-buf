# Configuration

Configuration is managed with Koanf, merging YAML files and environment variable overrides.

Files
- `config/local.yaml`: defaults for local dev
- `config/production.yaml`: sane defaults for prod; secret values are read from env

Resolution
- On startup, the app selects a file by `ENVIRONMENT`:
  - `dev`/empty: `config/local.yaml`
  - `prod`: `config/production.yaml`
- You can override file selection with `CONFIG_PATH=/path/to/config.yaml`.
- After loading, values are exported to environment variables so existing code continues to work.

Environment Overrides (Koanf)
- Prefix: `CFG_`
- Separator: `__` (double underscore) for nesting
- Examples:
  - `CFG_SERVER__PORT=9090` → `server.port`
  - `CFG_DATABASE__MAX_CONNS=200` → `database.max_conns`

Schema

```yaml
environment: dev|prod
server:
  port: 8080
  cors_allowed_origins: ["*"]
  run_migrations: true
  log_level: debug|info|warn|error
database:
  url: postgres://...
  max_conns: 50
  min_conns: 0
otel:
  endpoint: otel-collector:4317
  service_name: grpc-buf
security:
  jwt_secret: "..."
```

Validation
- In production:
  - `database.url` and `security.jwt_secret` are required.
- `server.port` must be 1-65535.

