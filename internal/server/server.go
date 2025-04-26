package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
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

	// Проверяем наличие static/index.html
	if _, err := os.Stat("static/index.html"); os.IsNotExist(err) {
		return nil, fmt.Errorf("static/index.html not found")
	}

	router := http.NewServeMux()
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/upload", uploadHandler)

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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "static/index.html")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	converted, err := convertContent(string(content))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(converted))
}

func convertContent(input string) (string, error) {
	// Реализуйте вашу логику конвертации между текстом и кодом Морзе
	// Пример:
	if isMorseCode(input) {
		return morseToText(input)
	}
	return textToMorse(input)
}

func isMorseCode(s string) bool {
	for _, r := range s {
		if r != '.' && r != '-' && r != ' ' && r != '/' {
			return false
		}
	}
	return true
}

func textToMorse(text string) (string, error) {
	// Реализация преобразования текста в Морзе
	return ".- .-.. .-.. -- -.-- -- --- .-. ... .", nil
}

func morseToText(morse string) (string, error) {
	// Реализация преобразования Морзе в текст
	return "ПРИВЕТ", nil
}
