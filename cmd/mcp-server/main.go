package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/grpc-buf/internal/config"
	"github.com/grpc-buf/internal/postgres"
	mcptransport "github.com/grpc-buf/internal/transport/mcp"
)

func main() {
	// Load configuration
	path, _ := config.ResolvePath()
	cfg, err := config.Load(path)
	if err == nil {
		if vErr := cfg.Validate(); vErr != nil {
			slog.Error("invalid configuration", "error", vErr)
			os.Exit(1)
		}
		// Configure slog level from config
		lvl := parseLevel(cfg.Server.LogLevel)
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl, AddSource: true}))
		slog.SetDefault(logger)
	} else {
		// In production, fail fast if config cannot be loaded
		env := strings.ToLower(strings.TrimSpace(os.Getenv("ENVIRONMENT")))
		if env == "prod" || env == "production" {
			slog.Error("failed to load configuration in production", "error", err)
			os.Exit(1)
		}
		// Dev: proceed with environment-only fallback
		slog.Warn("could not load config file; proceeding with environment only", "error", err)
		lvl := parseLevel(os.Getenv("LOG_LEVEL"))
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl, AddSource: true}))
		slog.SetDefault(logger)
	}

	slog.Info("Starting MCP server", "version", "1.0.0")

	// Initialize database connection
	dataStore := postgres.NewDatabaseConnectionFromConfig(cfg)
	defer dataStore.Close()

	// Create MCP server
	mcpServer, err := mcptransport.NewServer(dataStore)
	if err != nil {
		slog.Error("failed to create MCP server", "error", err)
		os.Exit(1)
	}

	// Setup signal handling for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start MCP server
	if err := mcpServer.Serve(ctx); err != nil {
		slog.Error("MCP server error", "error", err)
		os.Exit(1)
	}

	slog.Info("MCP server shutdown complete")
}

func parseLevel(v string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
