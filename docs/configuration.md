# Configuration

Configuration is managed with [`envconfig`](https://github.com/kelseyhightower/envconfig), layering environment variables on top of optional YAML files.

Files
- `config/local.yaml`: defaults for local dev
- `config/production.yaml`: sane defaults for prod; secret values are read from env

Resolution
- On startup, the app resolves a file path:
  - default: `config/local.yaml`
  - `CONFIG_PATH=/path/to/config.yaml` overrides the default

Environment Overrides
- Preferred: match the struct hierarchy directly, e.g. `SERVER_PORT`, `DATABASE_URL`, `SECURITY_JWT_SECRET`.
- Legacy: the `CFG_` prefix still works (`CFG_SERVER_PORT` maps to `server.port`).
- Examples:
  - `SERVER_PORT=9090` → `server.port`
  - `CFG_DATABASE_MAX_CONNS=200` → `database.max_conns`

Connection Pool Env Names
- Preferred: `DATABASE_MAX_CONNS`, `DATABASE_MIN_CONNS`.
- Legacy supported: `DB_MAX_CONNS`, `DB_MIN_CONNS`.

CORS Allowed Origins
- Dev: if `server.cors_allowed_origins` is empty, all origins are allowed.
- Any env: including `"*"` in `server.cors_allowed_origins` allows all origins.

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
