# Project Name: gRPC-buf

gRPC-buf is a Golang-based project that utilizes the `github.com/bufbuild/connect-go` library to develop gRPC and
REST APIs. It is designed to be deployed on Google Cloud Run and uses the `github.com/cloudevents/sdk-go` library for
publishing CloudEvents to Google Cloud Pub/Sub and Apache Kafka topics.

## Table of Contents

- [Project Name: gRPC-buf](#project-name-grpc-buf)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)

## Features

- gRPC and REST API development using `github.com/bufbuild/connect-go` library.
- Deployment on Google Cloud Run.
- CloudEvents integration with Google Cloud Pub/Sub and Apache Kafka topics using `github.com/cloudevents/sdk-go`.

## Prerequisites

- [Go](https://golang.org/doc/install) (1.17 or later)
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

## Installation

1. Clone the repository:

```bash
$ git clone https://github.com/dipjoytimetia/gRPC-buf.git
$ cd gRPC-buf
```