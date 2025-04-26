package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nikysoyd/sprint6/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "SERVER: ", log.LstdFlags|log.Lshortfile)

	// Создаем папку static если её нет
	if err := os.MkdirAll("static", 0755); err != nil {
		logger.Fatalf("Failed to create static dir: %v", err)
	}

	srv, err := server.NewServer(logger)
	if err != nil {
		logger.Fatalf("Server creation failed: %v", err)
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Println("Server starting on :8080")
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed: %v", err)
		}
	}()

	<-done
	logger.Println("Server stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server shutdown failed: %v", err)
	}
	logger.Println("Server stopped")
}
