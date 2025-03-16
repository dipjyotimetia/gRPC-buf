# gRPC-buf

A modern Golang service template featuring dual protocol APIs (gRPC/REST) using Connect, with deployment support for Google Cloud Run and event streaming capabilities.

[![Go Version](https://img.shields.io/github/go-mod/go-version/dipjyotimetia/gRPC-buf)](go.mod)
[![License](https://img.shields.io/github/license/dipjyotimetia/gRPC-buf)](LICENSE)

## Overview

gRPC-buf is a production-ready template for building cloud-native microservices with:

- 🚀 **Dual Protocol Support**: gRPC and REST endpoints via [connectrpc.com/connect](https://connectrpc.com)
- 🌩️ **Cloud-Native**: Ready for Google Cloud Run deployment
- 📨 **Event Streaming**: Built-in support for Pub/Sub and Kafka
- 🛠️ **Modern Tooling**: Streamlined development with Buf for Protocol Buffers
- 🔍 **Observability**: OpenTelemetry integration for monitoring
- 🗄️ **Persistence**: PostgreSQL database integration

## Features

### API Development
- Unified gRPC and REST API endpoints using Connect
- Automatic OpenAPI documentation generation
- Protocol Buffer validation and linting with Buf
- Type-safe API contracts

### Cloud Integration
- Containerized deployment to Google Cloud Run
- Google Pub/Sub integration for event-driven architectures
- Kafka support for messaging
- Comprehensive OpenTelemetry instrumentation
- PostgreSQL for reliable data persistence

### Developer Experience
- Fast development workflow with hot reload
- Comprehensive test suite with integration tests
- Makefile-based task automation
- Docker containerization for consistent environments
- Database migrations

## Prerequisites

- Go 1.23 or later
- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- [Buf CLI](https://docs.buf.build/installation)
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install) (for deployment)

## Quick Start

1. **Clone the repository**:
   ```bash
   git clone https://github.com/dipjyotimetia/gRPC-buf.git
   cd gRPC-buf
   ```

2. **Setup local environment**:
   ```bash
   make setup
   ```

3. **Generate code from Protocol Buffers**:
   ```bash
   make generate
   ```

4. **Start the service**:
   ```bash
   make run
   ```

5. **Access the API**:
   - REST API: http://localhost:8080
   - gRPC: localhost:8080

## Common Development Commands

| Command | Description |
|---------|-------------|
| `make test` | Run all tests |
| `make lint` | Run linters for Go code and Protocol Buffers |
| `make build` | Build the service binary |
| `make clean` | Clean up generated files |
| `make migrate-up` | Apply database migrations |
| `make migrate-down` | Roll back database migrations |

## Project Structure

```
.
├── cmd/                # Application entry points
│   ├── main.go         # Main application
│   └── server/         # Server implementation
├── internal/           # Private application code
│   ├── const/          # Constants
│   ├── gen/proto/      # Generated protocol buffer code
│   ├── logz/           # Logging utilities
│   ├── postgres/       # Database access layer
│   └── service/        # Business logic
├── proto/              # Protocol buffer definitions
│   ├── google/         # Google API definitions
│   ├── payment/        # Payment service definitions
│   └── registration/   # User registration definitions
└── scripts/            # Utility scripts
```

## Documentation

- [API Documentation](doc/about.md)
- [Architecture Overview](doc/about.md)
- [Deployment Guide](doc/about.md)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.
