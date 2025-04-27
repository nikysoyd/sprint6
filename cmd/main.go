package main

import (
	"log"
	"os"

	"github.com/nikysoyd/sprint6/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "morse-converter: ", log.LstdFlags|log.Lshortfile)

	srv := server.NewServer(logger)
	if err := srv.Start(); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
