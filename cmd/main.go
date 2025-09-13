package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/grpc-buf/cmd/server"
	"github.com/grpc-buf/internal/config"
)

func main() {
	// Load configuration and export to env for existing code paths
	path, _ := config.ResolvePath()
	cfg, err := config.Load(path)
	if err == nil {
		if vErr := cfg.Validate(); vErr != nil {
			slog.Error("invalid configuration", "error", vErr)
			os.Exit(1)
		}
		_ = cfg.ApplyEnv()
		// Configure slog level from config
		lvl := parseLevel(os.Getenv("LOG_LEVEL"))
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl, AddSource: true}))
		slog.SetDefault(logger)
	} else {
		// If config fails to load, continue with env only but warn
		slog.Warn("could not load config file; proceeding with environment only", "error", err)
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
