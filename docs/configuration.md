# Configuration

Configuration is managed with Koanf, merging YAML files and environment variable overrides.

Files
- `config/local.yaml`: defaults for local dev
- `config/production.yaml`: sane defaults for prod; secret values are read from env

Resolution
- On startup, the app resolves a file path:
  - default: `config/local.yaml`
  - `CONFIG_PATH=/path/to/config.yaml` overrides the default

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
security:
  jwt_secret: "..."
  jwt_issuer: grpc-buf
  jwt_audience: grpc-buf
  auth_skip_suffixes: ["/RegisterUser", "/LoginUser"]
```

Validation
- In production:
  - `database.url` and `security.jwt_secret` are required.
- `server.port` must be 1-65535.
