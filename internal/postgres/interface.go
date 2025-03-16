package postgres

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	"connectrpc.com/connect"
	paymentv1 "github.com/grpc-buf/internal/gen/proto/payment"
	userv1 "github.com/grpc-buf/internal/gen/proto/registration"
	"github.com/grpc-buf/internal/postgres/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // postgres driver
)

const ENV = "ENVIRONMENT"

type DataStore interface {
	MakePayment(ctx context.Context, req *connect.Request[paymentv1.PaymentRequest]) (*connect.Response[paymentv1.PaymentResponse], error)
	LoginUser(ctx context.Context, req *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error)
	RegisterUser(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error)
	Close()
}

// Store database connection
type Store struct {
	db *pgxpool.Pool
}

// NewDatabaseConnection creates a new PostgreSQL connection pool
func NewDatabaseConnection() DataStore {
	env := os.Getenv(ENV)
	var connectionString string

	if env == "dev" {
		slog.Info("Connecting to PostgreSQL local")
		connectionString = "postgres://postgres:postgres@postgres:5432/grpcbuf?sslmode=disable"
	} else {
		slog.Info("Connecting to PostgreSQL in production")
		connectionString = os.Getenv("DATABASE_URL")
	}

	// Configure connection pool
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		slog.Error("Unable to parse database connection string", "error", err)
		panic(err)
	}

	// Set pool configuration
	config.MaxConnIdleTime = 3 * time.Minute
	config.MaxConns = 300
	config.MinConns = 20

	// Create connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		panic(err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		slog.Error("Failed to ping database", "error", err)
		panic(err)
	}

	sqlDB, err := sql.Open("pgx", connectionString)
	if err != nil {
		slog.Error("Failed to open database for migrations", "error", err)
		panic(err)
	}
	defer sqlDB.Close()

	if err := migrations.RunMigrations(ctx, sqlDB); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		panic(err)
	}

	return &Store{
		db: pool,
	}
}

// Close closes the database connection pool
func (s *Store) Close() {
	if s.db != nil {
		s.db.Close()
	}
}
