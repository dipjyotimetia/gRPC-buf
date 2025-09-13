FROM golang:1.25 AS builder

ENV CGO_ENABLED=0
ENV GO111MODULE=on

WORKDIR /src

COPY go.* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

ARG VERSION=unknown
ARG COMMIT=unknown
ARG DATE=unknown

RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w -X github.com/grpc-buf/internal/version.Version=${VERSION} -X github.com/grpc-buf/internal/version.Commit=${COMMIT} -X github.com/grpc-buf/internal/version.Date=${DATE}" -o /out/server ./cmd/api
FROM gcr.io/distroless/base-debian12
WORKDIR /app
USER 65532:65532
COPY --from=builder /out/server /app/server
# Copy configs so container can run with CONFIG_PATH
COPY --from=builder /src/config /app/config
ENV CONFIG_PATH=/app/config/local.yaml
EXPOSE 8080
ENTRYPOINT ["/app/server"]
