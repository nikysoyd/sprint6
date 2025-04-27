package main

import (
	"log"
	"os"

	"github.com/nikysoyd/sprint6/internal/server"
)

func main() {

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	newSrv := server.NewServer(logger)

	logger.Println("Сервер запущен на http://localhost:8080")
	if err := newSrv.HTTP.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
