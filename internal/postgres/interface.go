package postgres

import (
    "context"
    "database/sql"
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
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DataStore interface {
	MakePayment(ctx context.Context, req *connect.Request[paymentv1.PaymentRequest]) (*connect.Response[paymentv1.PaymentResponse], error)
	LoginUser(ctx context.Context, req *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error)
	RegisterUser(ctx context.Context, req *connect.Request[userv1.RegisterRequest]) (*connect.Response[userv1.RegisterResponse], error)
	// Expense APIs
	CreateExpense(ctx context.Context, req *connect.Request[expensev1.CreateExpenseRequest]) (*connect.Response[expensev1.Expense], error)
	GetExpense(ctx context.Context, req *connect.Request[expensev1.GetExpenseRequest]) (*connect.Response[expensev1.Expense], error)
	ListExpenses(ctx context.Context, req *connect.Request[expensev1.ListExpensesRequest]) (*connect.Response[expensev1.ListExpensesResponse], error)
	UpdateExpense(ctx context.Context, req *connect.Request[expensev1.UpdateExpenseRequest]) (*connect.Response[expensev1.Expense], error)
	DeleteExpense(ctx context.Context, req *connect.Request[expensev1.DeleteExpenseRequest]) (*connect.Response[timestamppb.Timestamp], error)
	// Health
	Ping(ctx context.Context) error
	Close()
}

// Store database connection
type Store struct {
    db  *pgxpool.Pool
    sec config.SecurityConfig
}

// NewDatabaseConnection creates a new PostgreSQL connection pool (env-based for backward compatibility)
func NewDatabaseConnection() DataStore {
	cfg := &config.Config{
		Environment: os.Getenv("ENVIRONMENT"),
		Database: config.DatabaseConfig{
			URL:      os.Getenv("DATABASE_URL"),
			MaxConns: int(mustAtoi(os.Getenv("DB_MAX_CONNS"), 50)),
			MinConns: int(mustAtoi(os.Getenv("DB_MIN_CONNS"), 0)),
		},
		Server: config.ServerConfig{RunMigrations: os.Getenv("RUN_MIGRATIONS") != "false"},
	}
	return NewDatabaseConnectionFromConfig(cfg)
}

// NewDatabaseConnectionFromConfig creates a PostgreSQL pool using the provided config.
func NewDatabaseConnectionFromConfig(cfg *config.Config) DataStore {
    connectionString := cfg.Database.URL
    if strings.ToLower(cfg.Environment) == "dev" && connectionString == "" {
        slog.Info("Connecting to PostgreSQL local (dev)")
        connectionString = "postgres://postgres:postgres@postgres:5432/grpcbuf?sslmode=disable"
    } else {
        slog.Info("Connecting to PostgreSQL", "environment", cfg.Environment)
    }

    poolCfg, err := pgxpool.ParseConfig(connectionString)
    if err != nil {
        slog.Error("Unable to parse database connection string", "error", err)
        panic(err)
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

    // Create pool (does not necessarily connect immediately)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
    if err != nil {
        slog.Error("Failed to connect to database", "error", err)
        panic(err)
    }

    // Retry ping with backoff to allow database to become ready (e.g., in docker-compose)
    connectTimeout := 60 * time.Second
    if v := strings.TrimSpace(cfg.Database.ConnectTimeout); v != "" {
        if d, perr := time.ParseDuration(v); perr == nil && d > 0 {
            connectTimeout = d
        }
    }
    start := time.Now()
    var pingErr error
    backoff := 500 * time.Millisecond
    for {
        pctx, pcancel := context.WithTimeout(context.Background(), 3*time.Second)
        pingErr = pool.Ping(pctx)
        pcancel()
        if pingErr == nil {
            break
        }
        if time.Since(start) >= connectTimeout {
            slog.Error("Failed to ping database after timeout", "timeout", connectTimeout, "error", pingErr)
            panic(pingErr)
        }
        slog.Info("Database not ready, retrying", "error", pingErr, "backoff", backoff)
        time.Sleep(backoff)
        if backoff < 5*time.Second {
            backoff *= 2
        }
    }

    if cfg.Server.RunMigrations {
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
	} else {
		slog.Info("Skipping migrations as configured")
	}

    return &Store{db: pool, sec: cfg.Security}
}

func mustAtoi(s string, def int) int64 {
	if s == "" {
		return int64(def)
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return int64(def)
	}
	return n
}

// Close closes the database connection pool
func (s *Store) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

// Ping checks database connectivity
func (s *Store) Ping(ctx context.Context) error {
	if s.db == nil {
		return fmt.Errorf("db pool is nil")
	}
	return s.db.Ping(ctx)
}
