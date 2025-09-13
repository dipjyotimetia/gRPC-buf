package main

import (
    "context"
    "database/sql"
    "log/slog"
    "os"

    _ "github.com/jackc/pgx/v5/stdlib"

    "github.com/grpc-buf/internal/config"
    "github.com/grpc-buf/internal/postgres/migrations"
)

func main() {
    path, _ := config.ResolvePath()
    cfg, err := config.Load(path)
    if err != nil {
        slog.Error("failed to load config", "error", err)
        os.Exit(1)
    }
    if vErr := cfg.Validate(); vErr != nil {
        slog.Error("invalid configuration", "error", vErr)
        os.Exit(1)
    }
    dsn := cfg.Database.URL
    if dsn == "" {
        slog.Error("DATABASE_URL is required")
        os.Exit(1)
    }
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := migrations.RunMigrations(context.Background(), db); err != nil {
		slog.Error("migrations failed", "error", err)
		os.Exit(1)
	}
}
