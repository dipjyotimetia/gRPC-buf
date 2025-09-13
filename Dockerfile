FROM golang:1.25 AS builder

ENV CGO_ENABLED=0
ENV GO111MODULE=on

WORKDIR /src

COPY go.* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /out/server ./cmd
FROM gcr.io/distroless/base-debian12
WORKDIR /app
USER 65532:65532
COPY --from=builder /out/server /app/server
EXPOSE 8080
ENTRYPOINT ["/app/server"]
