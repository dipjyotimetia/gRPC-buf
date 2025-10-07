# Docker Compose Guide

## Overview
This project supports running services via docker-compose for local development and testing.

## Services

### 1. API Server (api)
- Port: 8080
- Dockerfile: `Dockerfile`
- Entrypoint: `/app/server` (from `cmd/api`)

### 2. MCP Server (mcp-server)
- Dockerfile: `Dockerfile.mcp`
- Entrypoint: `/app/mcp-server` (from `cmd/mcp-server`)
- Stdio-based transport for Model Context Protocol

### 3. PostgreSQL (postgres)
- Image: postgres:17.4
- Port: 5432
- Database: grpcbuf
- Credentials: postgres/postgres

### 4. pgAdmin (pgadmin)
- Profile: debug
- Port: 5050
- Only runs with `--profile debug`

## Quick Commands

### Start all services
```bash
make docker-up
# or
docker-compose up -d
```

### Start only MCP server
```bash
make docker-mcp-up
# or
docker-compose up -d postgres mcp-server
```

### View logs
```bash
# All services
make docker-logs

# MCP server only
make docker-mcp-logs
```

### Stop services
```bash
make docker-down
# or
docker-compose down
```

### Rebuild and restart
```bash
make docker-rebuild
# or
docker-compose down && docker-compose build && docker-compose up -d
```

### Restart MCP server
```bash
make docker-mcp-restart
# or
docker-compose restart mcp-server
```

## Configuration

### Environment Variables
Services use the configuration file at `/app/config/local.yaml` inside containers, which includes:
- Database URL: `postgres://postgres:postgres@postgres:5432/grpcbuf?sslmode=disable`
- Environment: `dev`

### Local Development
When running services locally (not in Docker):
- Use `.env` file for configuration
- Database URL: `postgres://postgres:postgres@localhost:5432/grpcbuf?sslmode=disable`
- Config path: `./config/local.yaml`

## Database Connection

The database hostname differs between docker-compose and local development:

| Environment | Hostname | Full URL |
|------------|----------|----------|
| Docker Compose | `postgres` | `postgres://postgres:postgres@postgres:5432/grpcbuf?sslmode=disable` |
| Local | `localhost` | `postgres://postgres:postgres@localhost:5432/grpcbuf?sslmode=disable` |

The `config/local.yaml` uses `postgres` (docker hostname), while `.env` uses `localhost` for local development.

## Troubleshooting

### MCP Server Connection Issues
If the MCP server can't connect to the database:
1. Check that postgres is healthy: `docker-compose ps`
2. View logs: `make docker-mcp-logs`
3. Verify config is correct: `docker-compose exec mcp-server cat /app/config/local.yaml`

### Port Conflicts
If port 5432 or 8080 is already in use:
1. Stop local PostgreSQL: `brew services stop postgresql`
2. Or change ports in `docker-compose.yaml`
