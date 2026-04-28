package main

import (
	"log/slog"
	"os"

	"github.com/grpc-buf/internal/config"
	"github.com/grpc-buf/internal/server"
)

func main() {
	cfg, err := config.Bootstrap()
	if err != nil {
		slog.Error("configuration error", "error", err)
		os.Exit(1)
	}

	if err := server.Run(cfg); err != nil {
		slog.Error("server run failed", "error", err)
		os.Exit(1)
	}
}
