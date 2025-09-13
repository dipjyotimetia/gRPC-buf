# Workflows

This document describes developer and CI/CD workflows for this repo.

Local Development
- Initial setup: `make setup` (installs buf + linters)
- Generate protos: `make generate` (buf generate)
- Run app: `make run` (uses config/local.yaml)
- Compose stack: `docker compose up -d` (Postgres, Jaeger, OTel, API)
- Lint: `make lint`
- Unit tests: `make test`
- Integration tests: start stack then `make test-integration`

Code Generation

```mermaid
flowchart LR
  P[.proto files]\nproto/* --> B[buf generate]\n
  B --> G1[internal/gen/proto/*]
  B --> G2[.../v1connect/*.connect.go]
```

Protos
- Linted with Buf’s DEFAULT rules.
- Breaking checks run in CI against `main`.
- Google API style: resource nouns and colon actions for RPCs.

CI
- Go CI: `.github/workflows/go-ci.yaml`
  - Go 1.25.x build, unit tests (-race), golangci-lint
  - Buf setup, lint, breaking (for changed proto files)
  - buf generate

- Integration Tests: `.github/workflows/integration-tests.yaml`
  - Triggered via `workflow_dispatch` or repo var `RUN_INTEGRATION_TESTS=true` on pushes to main
  - Runs `docker compose up -d --build`
  - Waits for `/livez`
  - Runs `go test -tags=integration ./tests/integration/... ./cmd/server -v`
  - Dumps container logs on failure; tears down stack

CD
- Cloud Run deploy: `.github/workflows/gcpdeploy.yaml`
  - Auth via workload identity / service account JSON
  - Builds and pushes image to GCR
  - Deploys to Cloud Run with HTTP/2 flag
  - Use Secret Manager for `DATABASE_URL` and set other envs as required

