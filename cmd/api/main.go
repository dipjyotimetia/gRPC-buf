package main

import (
    "log/slog"
    "os"
    "strings"

    "github.com/grpc-buf/internal/config"
    "github.com/grpc-buf/internal/server"
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

    if err := server.Run(cfg); err != nil {
        slog.Error("server run failed", "error", err)
    }
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
