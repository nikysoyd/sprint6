package main

import (
	"log"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	// Создаём логгер
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	// Создаём сервер
	srv := server.New(logger)

	// Запускаем сервер
	logger.Println("Сервер запущен на http://localhost:8080")
	if err := srv.HTTP.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
