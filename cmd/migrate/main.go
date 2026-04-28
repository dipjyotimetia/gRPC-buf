package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/grpc-buf/internal/config"
	"github.com/grpc-buf/internal/postgres/migrations"
)

const migrationTimeout = 2 * time.Minute

func main() {
	cfg, err := config.Bootstrap()
	if err != nil {
		slog.Error("configuration error", "error", err)
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

	ctx, cancel := context.WithTimeout(context.Background(), migrationTimeout)
	defer cancel()

	if err := migrations.RunMigrations(ctx, db); err != nil {
		slog.Error("migrations failed", "error", err)
		os.Exit(1)
	}
}
