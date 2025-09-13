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

migrate-up: ## Run migrations up
	migrate -path internal/postgres/migrations -database "$(DATABASE_URL)" up

migrate-down: ## Roll back migrations
	migrate -path internal/postgres/migrations -database "$(DATABASE_URL)" down

migrate-run: ## Run embedded migrations using Go binary
	$(GO) run ./cmd/migrate

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

.PHONY: run
run: ## Run the service locally
	ENVIRONMENT=dev $(GO) run ./cmd/api

.PHONY: setup
setup: ## Install dev tools
	$(MAKE) buf golangci-lint protoc-gen-go protoc-gen-go-grpc protoc-gen-connect-go
