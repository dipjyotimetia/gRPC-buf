package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var migrationFiles embed.FS

// RunMigrations runs database migrations
func RunMigrations(ctx context.Context, db *sql.DB) error {
	slog.Info("Running database migrations")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migrations driver: %w", err)
	}

	src, err := iofs.New(migrationFiles, ".")
	if err != nil {
		return fmt.Errorf("could not create migrations source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "grpcbuf", driver)
	if err != nil {
		return fmt.Errorf("could not create migrations instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	slog.Info("Database migrations completed successfully")
	return nil
}
