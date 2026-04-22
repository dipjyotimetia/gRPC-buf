package server

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"github.com/grpc-buf/internal/config"
	"github.com/grpc-buf/internal/postgres"
	"github.com/grpc-buf/internal/security"
	"github.com/grpc-buf/internal/service"
	httptransport "github.com/grpc-buf/internal/transport/http"
	authmw "github.com/grpc-buf/internal/transport/middleware/auth"
	"github.com/grpc-buf/internal/transport/middleware/ratelimit"
	"github.com/grpc-buf/internal/version"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const shutdownTimeout = 10 * time.Second

// Run starts the HTTP/gRPC server and blocks until SIGINT/SIGTERM is received
// or the underlying listener fails. cfg must be non-nil; callers should obtain
// it via config.Bootstrap.
func Run(cfg *config.Config) error {
	if cfg == nil {
		return errors.New("server.Run: cfg is nil")
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	db, err := postgres.NewDatabaseConnectionFromConfig(ctx, cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	paymentService := service.NewPaymentService(db)
	userService := service.NewUserService(db)
	expenseService := service.NewExpenseService(db)

	interceptors := buildInterceptors(cfg)

	mux := httptransport.NewMuxWithInterceptors(
		paymentService,
		userService,
		expenseService,
		interceptors...,
	)
	root := http.NewServeMux()
	root.Handle("/readyz", readyHandler(db))
	root.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(version.Get()); err != nil {
			slog.Warn("failed to encode /version response", "error", err)
		}
	})
	root.Handle("/", mux)

	srv := &http.Server{
		Addr: listenAddr(cfg),
		Handler: h2c.NewHandler(
			newCORS(cfg).Handler(root),
			&http2.Server{},
		),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}

	slog.Info("Starting gRPC server", "addr", srv.Addr)

	serveErr := make(chan error, 1)
	go func() {
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		// Buffered channel guarantees the send never blocks, even if nobody
		// reads (e.g. Shutdown errored and we returned early).
		serveErr <- err
	}()

	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	case err := <-serveErr:
		return err
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()
	shutdownErr := srv.Shutdown(shutdownCtx)
	if shutdownErr != nil {
		slog.Error("http shutdown failed", "error", shutdownErr)
	}
	if err := <-serveErr; err != nil {
		return err
	}
	return shutdownErr
}

func buildInterceptors(cfg *config.Config) []connect.Interceptor {
	loginRPS := 5
	loginBurst := 10
	if cfg.Server.LoginRPS > 0 {
		loginRPS = cfg.Server.LoginRPS
	}
	if cfg.Server.LoginBurst > 0 {
		loginBurst = cfg.Server.LoginBurst
	}
	interceptors := []connect.Interceptor{
		ratelimit.NewLoginInterceptor(float64(loginRPS), loginBurst),
	}

	verifier, err := security.NewVerifierFromConfig(cfg.Security)
	if err != nil || verifier == nil {
		slog.Warn("JWT auth disabled: missing secret")
		return interceptors
	}
	skip := cfg.Security.AuthSkipSuffixes
	if len(skip) == 0 {
		skip = []string{"/RegisterUser", "/LoginUser"}
	}
	return append(interceptors, authmw.NewJWTAuthInterceptor(verifier, skip))
}

// listenAddr resolves the bind address from (in order): Cloud Run's PORT,
// cfg.Server.Port, or the default :8080.
func listenAddr(cfg *config.Config) string {
	if p := strings.TrimSpace(os.Getenv("PORT")); p != "" {
		return ":" + p
	}
	if cfg.Server.Port != 0 {
		return ":" + strconv.Itoa(cfg.Server.Port)
	}
	return ":8080"
}

func newCORS(cfg *config.Config) *cors.Cors {
	origins := append([]string{}, cfg.Server.CORSAllowedOrigins...)
	env := strings.ToLower(strings.TrimSpace(cfg.Environment))
	allowAll := env == "dev" && len(origins) == 0
	if slices.Contains(origins, "*") {
		allowAll = true
	}

	return cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowOriginFunc: func(origin string) bool {
			if allowAll {
				return true
			}
			return slices.Contains(origins, origin)
		},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{
			"Accept",
			"Accept-Encoding",
			"Accept-Post",
			"Connect-Accept-Encoding",
			"Connect-Content-Encoding",
			"Content-Encoding",
			"Grpc-Accept-Encoding",
			"Grpc-Encoding",
			"Grpc-Message",
			"Grpc-Status",
			"Grpc-Status-Details-Bin",
		},
		MaxAge: int(2 * time.Hour / time.Second),
	})
}

// readyHandler returns HTTP 200 when dependencies are ready (DB ping).
func readyHandler(db postgres.DataStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := db.Ping(ctx); err != nil {
			http.Error(w, "not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("ok")); err != nil {
			slog.Debug("readyz write failed", "error", err)
		}
	})
}
