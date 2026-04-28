package migrations

import (
	"context"
	"database/sql"
	"embed"
	"errors"
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

	// m.Up() is blocking and ignores context, so run it in a goroutine and
	// signal GracefulStop when ctx fires. The buffered channel ensures the
	// goroutine can always exit even after we've stopped reading.
	done := make(chan error, 1)
	go func() { done <- m.Up() }()

	select {
	case err := <-done:
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("could not run migrations: %w", err)
		}
	case <-ctx.Done():
		// Ask migrate to stop at the next safe point. The goroutine will exit
		// once the in-flight statement completes; the process is expected to
		// exit shortly after this returns.
		select {
		case m.GracefulStop <- true:
		default:
		}
		return fmt.Errorf("migrations timed out: %w", ctx.Err())
	}

	slog.Info("Database migrations completed successfully")
	return nil
}
