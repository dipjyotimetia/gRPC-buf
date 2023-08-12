FROM golang:1.21.0 as builder

ENV CGO_ENABLED=0
ENV GO111MODULE=on

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

# Build the Go app
RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -a -o ./server ./cmd

# second stage
FROM debian:buster-slim
WORKDIR /app
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server /app/server

# Run the binary program produced by `go install`
ENTRYPOINT ["/app/server"]
