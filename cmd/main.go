package main

import (
	"log"

	"github.com/grpc-buf/cmd/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
