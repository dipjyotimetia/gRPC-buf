package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-buf/internal/config"
	"github.com/grpc-buf/internal/postgres"
	mcptransport "github.com/grpc-buf/internal/transport/mcp"
	"github.com/grpc-buf/internal/version"
)

func main() {
	cfg, err := config.Bootstrap()
	if err != nil {
		slog.Error("configuration error", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting MCP server", "version", version.Get().Version)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	dataStore, err := postgres.NewDatabaseConnectionFromConfig(ctx, cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dataStore.Close()

	mcpServer, err := mcptransport.NewServer(dataStore)
	if err != nil {
		slog.Error("failed to create MCP server", "error", err)
		os.Exit(1)
	}

	if err := mcpServer.Serve(ctx); err != nil {
		slog.Error("MCP server error", "error", err)
		os.Exit(1)
	}

	slog.Info("MCP server shutdown complete")
}
