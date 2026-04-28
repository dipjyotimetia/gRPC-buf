package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// ParseLogLevel maps a string level name to slog.Level. Unknown or empty
// values default to slog.LevelInfo so misconfiguration never silences logs.
func ParseLogLevel(v string) slog.Level {
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

// InstallLogger sets the process-wide slog.Default to a JSON handler at the
// given level with source location annotations.
func InstallLogger(level slog.Level) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}))
	slog.SetDefault(logger)
}

// Bootstrap resolves the config path, loads and validates the configuration,
// and installs the structured logger. In production it returns an error if the
// config file cannot be loaded. In dev it falls back to environment-only
// configuration and logs a warning.
//
// Callers should treat the returned *Config as authoritative: it is always
// non-nil on success.
func Bootstrap() (*Config, error) {
	path, err := ResolvePath()
	if err != nil {
		return nil, fmt.Errorf("resolve config path: %w", err)
	}

	cfg, loadErr := Load(path)
	if loadErr == nil {
		if vErr := cfg.Validate(); vErr != nil {
			return nil, fmt.Errorf("invalid configuration: %w", vErr)
		}
		InstallLogger(ParseLogLevel(cfg.Server.LogLevel))
		return cfg, nil
	}

	env := strings.ToLower(strings.TrimSpace(os.Getenv("ENVIRONMENT")))
	if env == "prod" || env == "production" {
		return nil, fmt.Errorf("failed to load configuration in production: %w", loadErr)
	}

	// Dev fallback: configure logging from LOG_LEVEL and load env-only config.
	InstallLogger(ParseLogLevel(os.Getenv("LOG_LEVEL")))
	slog.Warn("could not load config file; proceeding with environment only", "error", loadErr)

	envCfg, envErr := Load("")
	if envErr != nil {
		return nil, errors.Join(loadErr, envErr)
	}
	return envCfg, nil
}
