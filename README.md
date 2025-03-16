# gRPC-buf

A modern Golang service template featuring gRPC/REST APIs using Connect, with deployment support for Google Cloud Run and event streaming capabilities.

[![Go Version](https://img.shields.io/github/go-mod/go-version/dipjyotimetia/gRPC-buf)](go.mod)
[![License](https://img.shields.io/github/license/dipjyotimetia/gRPC-buf)](LICENSE)

## Overview

gRPC-buf provides a production-ready template for building microservices with:

- 🚀 Dual protocol support (gRPC and REST) using [connectrpc.com/connect](https://connectrpc.com)
- 🌩️ Cloud-native deployment on Google Cloud Run
- 📨 Event streaming with CloudEvents to Pub/Sub and Kafka
- 🛠️ Modern development workflow with Buf for Protocol Buffers
- 🗄️ PostgreSQL for persistent data storage

## Features

- **API Development**
  - gRPC and REST endpoints using Connect
  - Automatic OpenAPI documentation
  - Protocol Buffer validation and linting
  
- **Cloud Integration**
  - Google Cloud Run deployment
  - Pub/Sub integration
  - Kafka support
  - OpenTelemetry observability
  - PostgreSQL for data persistence

- **Developer Experience**
  - Hot reload during development
  - Integrated testing framework
  - Makefile automation
  - Docker containerization

## Prerequisites

- Go 1.23 or later
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Buf CLI](https://docs.buf.build/installation)

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/dipjyotimetia/gRPC-buf.git
cd gRPC-buf
```

2. Install dependencies:
```bash
make setup
```

3. Generate protobuf code:
```bash
make generate
```

4. Run the service:
```bash
make run
```

## Development

```bash
# Run tests
make test

# Lint code
make lint

# Build binary
make build

# Clean up
make clean
```

## Documentation

- API Documentation
- Architecture Overview
- Deployment Guide

## Contributing

Contributions are welcome! Please read our Contributing Guide for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License.
