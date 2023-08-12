package main

import (
	"log/slog"

	"github.com/grpc-buf/cmd/server"
)

func main() {
	if err := server.Run(); err != nil {
		slog.Error(err.Error())
	}
}
