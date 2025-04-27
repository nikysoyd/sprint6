package server

import (
	"log"
	"net/http"
	"time"

	"github.com/nikysoyd/sprint6/internal/handlers"
)

type Server struct {
	Logger *log.Logger
	HTTP   *http.Server
}

func NewServer(logger *log.Logger) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.RootServeHome)
	mux.HandleFunc("/upload", handlers.UploadFileHandler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		Logger: logger,
		HTTP:   srv,
	}
}
