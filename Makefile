# See https://tech.davis-hansson.com/p/make/
SHELL := bash
.DELETE_ON_ERROR:
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-print-directory
COPYRIGHT_YEARS := 2022
LICENSE_IGNORE := -e /testdata/
# Set to use a different compiler. For example, `GO=go1.18rc1 make test`.
GO ?= go
# Load environment variables from .env if present (for local dev)
ifneq (,$(wildcard .env))
include .env
export $(shell sed -n 's/^[[:space:]]*\([A-Za-z_][A-Za-z0-9_]*\)[[:space:]]*=.*/\1/p' .env)
endif
# Optional DSN override for migration targets
DSN ?= $(DATABASE_URL)

.PHONY: help
help: ## Describe useful make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-30s %s\n", $$1, $$2}'

.PHONY: all
all: ## Build, test, and lint (default)
	$(MAKE) test
	$(MAKE) lint

.PHONY: clean
clean: ## Delete intermediate build artifacts
	@# -X only removes untracked files, -d recurses into directories, -f actually removes files/dirs
	git clean -Xdf

.PHONY: test
test: build ## Run unit tests
	$(GO) test -vet=off -race -cover ./...

.PHONY: test-integration
test-integration: ## Run integration tests (server must be running)
	$(GO) test -tags=integration ./tests/integration/... ./internal/server -v

.PHONY: build
build: generate ## Build all packages
	$(GO) build ./...

.PHONY: lint
lint: golangci-lint buf ## Lint Go and protobuf
	test -z "$$(buf format -d . | tee /dev/stderr)"
	$(GO) vet ./...
	golangci-lint run
	buf lint

.PHONY: lintfix
lintfix: golangci-lint buf ## Automatically fix some lint errors
	golangci-lint run --fix
	buf format -w .

.PHONY: upgrade
upgrade: ## Upgrade dependencies
	go get -u -t ./... && go mod tidy -v

# Migration commands
.PHONY: migrate-create migrate-up migrate-down migrate-run

migrate-create: ## Create a new migration file
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir internal/postgres/migrations -seq $$name

migrate-up: ## Run migrations up (external tool). Set DSN=postgres://...
	@dsn="$(DSN)"; \
	case "$$dsn" in \
		postgresql://*) echo "[migrate-up] Normalizing DSN scheme postgresql:// -> postgres://"; dsn="$${dsn#postgresql://}"; dsn="postgres://$${dsn}";; \
	esac; \
	migrate -path internal/postgres/migrations -database "$$dsn" up

migrate-down: ## Roll back migrations (external tool). Set DSN=postgres://...
	@dsn="$(DSN)"; \
	case "$$dsn" in \
		postgresql://*) echo "[migrate-down] Normalizing DSN scheme postgresql:// -> postgres://"; dsn="$${dsn#postgresql://}"; dsn="postgres://$${dsn}";; \
	esac; \
	migrate -path internal/postgres/migrations -database "$$dsn" down

migrate-run: ## Run embedded migrations using Go binary (uses config)
	$(GO) run ./cmd/migrate

.PHONY: migrate-run-local
migrate-run-local: ## Run embedded migrations against custom DSN. Usage: make migrate-run-local DSN=postgres://...
	@dsn="$(DSN)"; \
	case "$$dsn" in \
		postgresql://*) echo "[migrate-run-local] Normalizing DSN scheme postgresql:// -> postgres://"; dsn="$${dsn#postgresql://}"; dsn="postgres://$${dsn}";; \
	esac; \
	CONFIG_PATH=config/local.yaml CFG_DATABASE_URL="$$dsn" $(GO) run ./cmd/migrate

.PHONY: buf golangci-lint protoc-gen-go protoc-gen-go-grpc

buf:
	go install github.com/bufbuild/buf/cmd/buf@latest

golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

protoc-gen-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

protoc-gen-go-grpc:
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: protoc-gen-connect-go
protoc-gen-connect-go:
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest

.PHONY: generate
generate: ## Generate code from protobufs
	buf generate
	@echo "Removing generated google API files (using genproto instead)..."
	rm -rf internal/gen/proto/google

.PHONY: run
run: ## Run the service locally
	ENVIRONMENT=dev $(GO) run ./cmd/api

.PHONY: mcp-build
mcp-build: generate ## Build MCP server binary
	$(GO) build -o bin/mcp-server ./cmd/mcp-server

.PHONY: mcp-run
mcp-run: ## Run MCP server locally
	ENVIRONMENT=dev $(GO) run ./cmd/mcp-server

.PHONY: mcp-generate
mcp-generate: ## Generate MCP stubs from protobuf
	buf generate
	@echo "Removing generated google API files (using genproto instead)..."
	rm -rf internal/gen/proto/google

.PHONY: docker-up
docker-up: ## Start all services with docker-compose
	docker-compose up -d

.PHONY: docker-down
docker-down: ## Stop all docker-compose services
	docker-compose down

.PHONY: docker-logs
docker-logs: ## View logs from all docker-compose services
	docker-compose logs -f

.PHONY: docker-mcp-logs
docker-mcp-logs: ## View logs from MCP server
	docker-compose logs -f mcp-server

.PHONY: docker-build
docker-build: ## Build docker images
	docker-compose build

.PHONY: docker-rebuild
docker-rebuild: ## Rebuild and restart all services
	docker-compose down
	docker-compose build
	docker-compose up -d

.PHONY: docker-mcp-up
docker-mcp-up: ## Start only MCP server and dependencies
	docker-compose up -d postgres mcp-server

.PHONY: docker-mcp-restart
docker-mcp-restart: ## Restart MCP server
	docker-compose restart mcp-server

.PHONY: setup
setup: ## Install dev tools
	$(MAKE) buf golangci-lint protoc-gen-go protoc-gen-go-grpc protoc-gen-connect-go
.PHONY: migrate-install
migrate-install: ## Install golang-migrate CLI with Postgres driver (requires Go toolchain)
	CGO_ENABLED=0 $(GO) install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Installed migrate CLI with Postgres driver. Ensure \"$$(go env GOPATH)/bin\" is on your PATH."

.PHONY: migrate-docker
migrate-docker: ## Run migrations using migrate/migrate Docker image. Usage: make migrate-docker DSN=postgres://...
	@dsn="$(DSN)"; \
	case "$$dsn" in \
		postgresql://*) echo "[migrate-docker] Normalizing DSN scheme postgresql:// -> postgres://"; dsn="$${dsn#postgresql://}"; dsn="postgres://$${dsn}";; \
	esac; \
	docker run --rm -v "$$PWD/internal/postgres/migrations:/migrations" migrate/migrate:latest -path=/migrations -database "$$dsn" up
