package main

import (
	"log"
	"net/http"
	"os"

	"github.com/nikysoyd/sprint6/internal/server"
)

func main() {
	// Создаём логгер
	logger := log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lshortfile)

	// Создаём сервер (теперь с обработкой ошибки)
	srv, err := server.NewServer(logger)
	if err != nil {
		logger.Fatalf("Failed to create server: %v", err)
	}

	// Запускаем сервер
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
