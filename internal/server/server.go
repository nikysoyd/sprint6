package server

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/nikysoyd/sprint6/pkg/morse" // Импортируем ваш пакет morse
)

type Server struct {
	logger *log.Logger
	server *http.Server
}

func NewServer(logger *log.Logger) (*Server, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("<html><body><h1>Морзе-конвертер</h1></body></html>"))
	})

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

		text := string(body)
		morseCode := morse.ToMorse(text)
		w.Write([]byte(morseCode))
	})

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

		morseCode := string(body)
		text := morse.ToText(morseCode)
		w.Write([]byte(text))
	})

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

		content := string(body)

		// Определяем тип контента (текст или морзе)
		if strings.ContainsAny(content, ".-") {
			// Это код Морзе
			text := morse.ToText(content)
			w.Write([]byte(text))
		} else {
			// Это обычный текст
			morseCode := morse.ToMorse(content)
			w.Write([]byte(morseCode))
		}
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
