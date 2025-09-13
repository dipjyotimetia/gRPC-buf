package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/grpc-buf/internal/config"
	constant "github.com/grpc-buf/internal/const"
	"github.com/grpc-buf/internal/logz"
	"github.com/grpc-buf/internal/postgres"
	"github.com/grpc-buf/internal/service"
	"github.com/rs/cors"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func Run(cfg *config.Config) error {
	// Initialize dependencies at runtime using config
	var db postgres.DataStore
	if cfg != nil {
		db = postgres.NewDatabaseConnectionFromConfig(cfg)
	} else {
		db = postgres.NewDatabaseConnection()
	}
	defer db.Close()

	paymentService := service.NewPaymentService(db)
	userService := service.NewUserService(db)
	expenseService := service.NewExpenseService(db)

	mux := setupHandler(paymentService, userService, expenseService)
	// Configure slog: JSON to stdout, include source locations
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	slog.SetDefault(log)
	log.Info("Starting gRPC server")

	otelShutdown, err := logz.StartTracer()
	if err != nil {
		log.Error("error setting up OTel SDK", "error", err)
	}
	defer otelShutdown()

	ctx, span := constant.Tracer.Start(context.Background(), "server_startup")
	defer span.End()

	addr := ":8080"
	if cfg != nil && cfg.Server.Port != 0 {
		addr = ":" + os.Getenv("PORT")
	} else if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	srv := &http.Server{
		Addr: addr,
		Handler: h2c.NewHandler(
			newCORS().Handler(mux),
			&http2.Server{},
		),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http listen and serve failed", "error", err)
		}
	}()

	<-signals

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "server shutdown failed")
		log.Error("http shutdown failed", "error", err)
	}
	return nil
}

func newCORS() *cors.Cors {
	// Configurable CORS via env; default permissive in dev only
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	origins := []string{}
	if allowedOrigins != "" {
		for _, o := range strings.Split(allowedOrigins, ",") {
			if trimmed := strings.TrimSpace(o); trimmed != "" {
				origins = append(origins, trimmed)
			}
		}
	}

	allowAll := os.Getenv("ENVIRONMENT") == "dev" && len(origins) == 0

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
