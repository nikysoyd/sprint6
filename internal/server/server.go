package server

import (
	"io"
	"log"
	"net/http"

	"github.com/nikysoyd/sprint6/pkg/morse"
)

type Server struct {
	logger *log.Logger
	server *http.Server
}

func NewServer(logger *log.Logger) (*Server, error) {
	mux := http.NewServeMux()

	// Главная страница
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("<html><body><h1>Конвертер Морзе</h1></body></html>"))
	})

	// Текст в Морзе
	mux.HandleFunc("/text-to-morse", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		result := morse.ToMorse(string(body))
		w.Write([]byte(result))
	})

	// Морзе в текст
	mux.HandleFunc("/morse-to-text", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		result := morse.ToText(string(body))
		w.Write([]byte(result))
	})

	// Автоматическое определение
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		input := string(body)
		var result string

		// Проверяем, является ли ввод кодом Морзе
		isMorse := true
		for _, r := range input {
			if r != '.' && r != '-' && r != ' ' && r != '/' {
				isMorse = false
				break
			}
		}

		if isMorse {
			result = morse.ToText(input)
		} else {
			result = morse.ToMorse(input)
		}

		w.Write([]byte(result))
	})

	server := &http.Server{
		Addr:     ":8080",
		Handler:  mux,
		ErrorLog: logger,
	}

	return &Server{logger: logger, server: server}, nil
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}
