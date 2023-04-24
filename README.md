# Project Name: gRPC-buf

gRPC-buf is a Golang-based project that utilizes the `github.com/bufbuild/connect-go` library to develop gRPC and
REST APIs. It is designed to be deployed on Google Cloud Run and uses the `github.com/cloudevents/sdk-go` library for
publishing CloudEvents to Google Cloud Pub/Sub and Apache Kafka topics.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Deployment](#deployment)

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
$ git clone https://github.com/yourusername/GoGRPCRestAPI.git
$ cd gRPC-buf
```

2. Deploy the Cloud Run service:

```bash

$ docker build -t gcr.io/your-project-id/go-grpc-rest-api:v1 .
$ gcloud auth configure-docker
$ docker push gcr.io/your-project-id/go-grpc-rest-api:v1
$ gcloud run deploy go-grpc-rest-api --image gcr.io/your-project-id/go-grpc-rest-api:v1 --platform managed --region your-region --allow-unauthenticated
$ SERVICE_URL=$(gcloud run services describe go-grpc-rest-api --platform managed --region your-region --format 'value(status.url)')
$ curl $SERVICE_URL

```