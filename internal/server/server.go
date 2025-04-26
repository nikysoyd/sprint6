package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	listener   net.Listener
	logger     *log.Logger
}

func NewServer(logger *log.Logger) (*Server, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	router := http.NewServeMux()
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/upload", uploadHandler)

	// Создаем listener заранее, чтобы проверить доступность порта
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	return &Server{
		httpServer: &http.Server{
			Handler:      router,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
		listener: listener,
		logger:   logger,
	}, nil
}

func (s *Server) Start() error {
	s.logger.Printf("Server starting on %s", s.listener.Addr())
	return s.httpServer.Serve(s.listener)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "static/index.html")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Обработка загрузки файла...
}
