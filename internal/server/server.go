package server

import (
	"log"
	"net/http"
)

type Server struct {
	logger *log.Logger
	server *http.Server
}

func NewServer(logger *log.Logger) (*Server, error) {
	mux := http.NewServeMux()
	// Пример: регистрируем хендлеры
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	server := &http.Server{
		Addr:     ":8080",
		Handler:  mux,
		ErrorLog: logger,
	}

	return &Server{
		logger: logger,
		server: server,
	}, nil
}
