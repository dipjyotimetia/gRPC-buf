package main

import (
	"log"

	"github.com/grpc-buf/cmd/server"
	// This import path is based on the name declaration in the go.mod,
	// and the gen/proto/go output location in the buf.gen.yaml.
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
