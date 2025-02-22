package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	constant "github.com/grpc-buf/internal/const"
	"github.com/grpc-buf/internal/logz"
	"github.com/grpc-buf/internal/mongo"
	"github.com/grpc-buf/internal/service"
	"github.com/rs/cors"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	db             = mongo.NewDatabaseConnection()
	paymentService = service.NewPaymentService(db)
	userService    = service.NewUserService(db)
)

func Run() error {
	mux := setupHandler()
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	log.Info("Starting gRPC server")

	otelShutdown, err := logz.StartTracer()
	if err != nil {
		log.Error("error setting up OTel SDK - %e")
	}
	defer otelShutdown()

	ctx, span := constant.Tracer.Start(context.Background(), "server_startup")
	defer span.End()

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
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
			log.Error("HTTP listen and serve: %v", err)
		}
	}()

	<-signals

	shutdownCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "server shutdown failed")
		log.Error("HTTP shutdown failed", "error", err)
	}
	return nil
}

func newCORS() *cors.Cors {
	// To let web developers play with the demo service from browsers, we need a
	// very permissive CORS setup.
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
			// Allow all origins, which effectively disables CORS.
			return true
		},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{
			// Content-Type is in the default safelist.
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
		// Let browsers cache CORS information for longer, which reduces the number
		// of preflight requests. Any changes to ExposedHeaders won't take effect
		// until the cached data expires. FF caps this value at 24h, and modern
		// Chrome caps it at 2h.
		MaxAge: int(2 * time.Hour / time.Second),
	})
}
