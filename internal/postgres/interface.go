package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/grpc-buf/internal/config"
	expensev1 "github.com/grpc-buf/internal/gen/proto/expense"
	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	userv1 "github.com/grpc-buf/internal/gen/proto/registration"
	"github.com/grpc-buf/internal/postgres/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DataStore is the persistence-facing contract used by both the Connect
// services and the MCP transport.
type DataStore interface {
	MakePayment(ctx context.Context, req *connect.Request[paymentv1.PaymentRequest]) (*connect.Response[paymentv1.PaymentResponse], error)
	LoginUser(ctx context.Context, req *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error)
	RegisterUser(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error)
	// Expense APIs
	CreateExpense(ctx context.Context, req *connect.Request[expensev1.CreateExpenseRequest]) (*connect.Response[expensev1.Expense], error)
	GetExpense(ctx context.Context, req *connect.Request[expensev1.GetExpenseRequest]) (*connect.Response[expensev1.Expense], error)
	ListExpenses(ctx context.Context, req *connect.Request[expensev1.ListExpensesRequest]) (*connect.Response[expensev1.ListExpensesResponse], error)
	UpdateExpense(ctx context.Context, req *connect.Request[expensev1.UpdateExpenseRequest]) (*connect.Response[expensev1.Expense], error)
	DeleteExpense(ctx context.Context, req *connect.Request[expensev1.DeleteExpenseRequest]) (*connect.Response[emptypb.Empty], error)
	// Health
	Ping(ctx context.Context) error
	Close()
}

// Store is the pgx-backed implementation of DataStore.
type Store struct {
	db  *pgxpool.Pool
	sec config.SecurityConfig
}

// NewDatabaseConnection creates a PostgreSQL pool from process environment
// variables. It is retained for call sites that predate config-driven wiring.
func NewDatabaseConnection(ctx context.Context) (DataStore, error) {
	maxConns, err := atoiOrDefault(envOr("DATABASE_MAX_CONNS", "DB_MAX_CONNS"), 50)
	if err != nil {
		return nil, fmt.Errorf("parse max conns: %w", err)
	}
	minConns, err := atoiOrDefault(envOr("DATABASE_MIN_CONNS", "DB_MIN_CONNS"), 0)
	if err != nil {
		return nil, fmt.Errorf("parse min conns: %w", err)
	}
	cfg := &config.Config{
		Environment: os.Getenv("ENVIRONMENT"),
		Database: config.DatabaseConfig{
			URL:      os.Getenv("DATABASE_URL"),
			MaxConns: maxConns,
			MinConns: minConns,
		},
		Server: config.ServerConfig{RunMigrations: os.Getenv("RUN_MIGRATIONS") != "false"},
	}
	return NewDatabaseConnectionFromConfig(ctx, cfg)
}

// NewDatabaseConnectionFromConfig creates a pgx connection pool, waits until
// the database is reachable (bounded by cfg.Database.ConnectTimeout or 60s),
// and optionally runs migrations. It returns an error instead of panicking so
// callers can report and exit gracefully.
func NewDatabaseConnectionFromConfig(ctx context.Context, cfg *config.Config) (DataStore, error) {
	connectionString := cfg.Database.URL
	if strings.ToLower(cfg.Environment) == "dev" && connectionString == "" {
		slog.Info("Connecting to PostgreSQL local (dev)")
		connectionString = "postgres://postgres:postgres@postgres:5432/grpcbuf?sslmode=disable"
	} else {
		slog.Info("Connecting to PostgreSQL", "environment", cfg.Environment)
	}

	poolCfg, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse database connection string: %w", err)
	}
	poolCfg.MaxConnIdleTime = 3 * time.Minute
	if cfg.Database.MaxConns > 0 {
		poolCfg.MaxConns = int32(cfg.Database.MaxConns)
	} else {
		poolCfg.MaxConns = 50
	}
	if cfg.Database.MinConns >= 0 {
		poolCfg.MinConns = int32(cfg.Database.MinConns)
	}

	poolCtx, poolCancel := context.WithTimeout(ctx, 10*time.Second)
	defer poolCancel()
	pool, err := pgxpool.NewWithConfig(poolCtx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if err := waitForDatabase(ctx, pool, cfg.Database.ConnectTimeout); err != nil {
		pool.Close()
		return nil, err
	}

	if cfg.Server.RunMigrations {
		if err := runMigrations(ctx, connectionString); err != nil {
			pool.Close()
			return nil, err
		}
	} else {
		slog.Info("Skipping migrations as configured")
	}

	return &Store{db: pool, sec: cfg.Security}, nil
}

// waitForDatabase pings the pool with exponential backoff until it succeeds or
// the timeout elapses. It honors ctx cancellation so a shutdown signal during
// startup stops the retry loop immediately.
func waitForDatabase(ctx context.Context, pool *pgxpool.Pool, timeoutStr string) error {
	connectTimeout := 60 * time.Second
	if v := strings.TrimSpace(timeoutStr); v != "" {
		if d, perr := time.ParseDuration(v); perr == nil && d > 0 {
			connectTimeout = d
		}
	}

	deadline := time.Now().Add(connectTimeout)
	backoff := 500 * time.Millisecond
	for {
		pctx, pcancel := context.WithTimeout(ctx, 3*time.Second)
		err := pool.Ping(pctx)
		pcancel()
		if err == nil {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("ping database after %s: %w", connectTimeout, err)
		}

		slog.Info("Database not ready, retrying", "error", err, "backoff", backoff)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
		if backoff < 5*time.Second {
			backoff *= 2
		}
	}
}

func runMigrations(ctx context.Context, connectionString string) error {
	sqlDB, err := sql.Open("pgx", connectionString)
	if err != nil {
		return fmt.Errorf("open database for migrations: %w", err)
	}
	defer sqlDB.Close()
	if err := migrations.RunMigrations(ctx, sqlDB); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	return nil
}

// atoiOrDefault parses s as an int. Empty strings return def; any other parse
// failure surfaces as an error so config mistakes are never silent.
func atoiOrDefault(s string, def int) (int, error) {
	if s == "" {
		return def, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("parse %q: %w", s, err)
	}
	return n, nil
}

func envOr(keys ...string) string {
	for _, k := range keys {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

// Close closes the underlying connection pool. Safe to call on a nil pool.
func (s *Store) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

// Ping checks database connectivity.
func (s *Store) Ping(ctx context.Context) error {
	if s.db == nil {
		return errors.New("db pool is nil")
	}
	return s.db.Ping(ctx)
}
