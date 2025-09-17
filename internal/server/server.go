package server

import (
    "context"
    "encoding/json"
    "errors"
    "github.com/grpc-buf/internal/transport/middleware/ratelimit"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "slices"
    "strings"
    "syscall"
    "time"

    "github.com/grpc-buf/internal/config"
    "github.com/grpc-buf/internal/postgres"
    "github.com/grpc-buf/internal/service"
    httptransport "github.com/grpc-buf/internal/transport/http"
    authmw "github.com/grpc-buf/internal/transport/middleware/auth"
    "connectrpc.com/connect"
    "github.com/grpc-buf/internal/version"
    "github.com/grpc-buf/internal/security"
    "github.com/rs/cors"
    "golang.org/x/net/http2"
    "golang.org/x/net/http2/h2c"
    "strconv"
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

    // Install middleware (rate limiting) from config
    loginRPS := 5
    loginBurst := 10
    if cfg != nil {
        if cfg.Server.LoginRPS > 0 {
            loginRPS = cfg.Server.LoginRPS
        }
        if cfg.Server.LoginBurst > 0 {
            loginBurst = cfg.Server.LoginBurst
        }
    }
    limiter := ratelimit.NewLoginInterceptor(float64(loginRPS), loginBurst)
    interceptors := []connect.Interceptor{limiter}
    // Auth from config (preferred), else env fallback
    var verifier *security.Verifier
    var vErr error
    if cfg != nil {
        verifier, vErr = security.NewVerifierFromConfig(cfg.Security)
    } else {
        verifier, vErr = security.NewVerifierFromEnv()
    }
    if vErr == nil && verifier != nil {
        skip := []string{"/RegisterUser", "/LoginUser"}
        if cfg != nil && len(cfg.Security.AuthSkipSuffixes) > 0 {
            skip = cfg.Security.AuthSkipSuffixes
        }
        jwtAuth := authmw.NewJWTAuthInterceptor(verifier, skip)
        interceptors = append(interceptors, jwtAuth)
    } else {
        slog.Warn("JWT auth disabled: missing secret")
    }

    mux := httptransport.NewMuxWithInterceptors(
        paymentService,
        userService,
        expenseService,
        interceptors...,
    )
    // Wrap with a root mux to add readiness and version endpoints
    root := http.NewServeMux()
    root.Handle("/readyz", readyHandler(db))
    root.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(version.Get())
    })
    root.Handle("/", mux)
    // Use preconfigured global logger (set in main) and log startup
    slog.Info("Starting gRPC server")

	ctx := context.Background()

    // Determine address: prefer Cloud Run $PORT env, else config, else default 8080
    addr := ":8080"
    if p := strings.TrimSpace(os.Getenv("PORT")); p != "" { // Cloud Run sets PORT
        addr = ":" + p
    } else if cfg != nil && cfg.Server.Port != 0 {
        addr = ":" + strconv.Itoa(cfg.Server.Port)
    }
    srv := &http.Server{
        Addr: addr,
        Handler: h2c.NewHandler(
            newCORS(cfg).Handler(root),
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
            slog.Error("http listen and serve failed", "error", err)
        }
    }()

	<-signals

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

    if err := srv.Shutdown(shutdownCtx); err != nil {
        slog.Error("http shutdown failed", "error", err)
    }
    return nil
}

func newCORS(cfg *config.Config) *cors.Cors {
    // Use config-driven allowed origins. In dev, allow all if unset.
    origins := []string{}
    if cfg != nil && len(cfg.Server.CORSAllowedOrigins) > 0 {
        origins = append(origins, cfg.Server.CORSAllowedOrigins...)
    }
    env := ""
    if cfg != nil { env = strings.ToLower(strings.TrimSpace(cfg.Environment)) }
    allowAll := env == "dev" && len(origins) == 0
    // Treat literal "*" as allow-all in any env.
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
		_, _ = w.Write([]byte("ok"))
	})
}

//
