package main

import (
	"log"
	"net/http"
	"os"

	"github.com/nikysoyd/sprint6/server" // Замените на актуальный путь к пакету server
)

func main() {
	// Создаём логгер
	logger := log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lshortfile)

	// Создаём сервер
	srv := server.NewServer(logger)

	// Запускаем сервер и обрабатываем ошибки
	if err := srv.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
