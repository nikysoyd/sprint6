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

// Server представляет HTTP сервер приложения
type Server struct {
	httpServer *http.Server
	listener   net.Listener
	logger     *log.Logger
}

// NewServer создает новый экземпляр сервера
func NewServer(logger *log.Logger) (*Server, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	// Проверяем наличие index.html в корне, создаем если нет
	if err := ensureIndexHTML(); err != nil {
		return nil, fmt.Errorf("failed to ensure index.html: %w", err)
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

// ensureIndexHTML проверяет наличие index.html и создает если нужно
func ensureIndexHTML() error {
	const filename = "index.html"

	// Проверяем существует ли файл
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Создаем файл с минимальным содержимым
		content := `<!DOCTYPE html>
<html>
<head>
    <title>Конвертер Морзе</title>
</head>
<body>
    <h1>Конвертер текста и кода Морзе</h1>
    <form action="/upload" method="post" enctype="multipart/form-data">
        <input type="file" name="file" required>
        <button type="submit">Конвертировать</button>
    </form>
</body>
</html>`

		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create %s: %w", filename, err)
		}
	}
	return nil
}

// Start запускает сервер
func (s *Server) Start() error {
	s.logger.Printf("Server starting on %s", s.listener.Addr())
	return s.httpServer.Serve(s.listener)
}

// Shutdown корректно останавливает сервер
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// rootHandler обрабатывает запросы к корневому пути
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "index.html")
}

// uploadHandler обрабатывает загрузку файлов
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсим multipart форму (макс. 10MB файл)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Получаем файл из формы
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Читаем содержимое файла
	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file content", http.StatusInternalServerError)
		return
	}

	// Конвертируем содержимое
	converted, err := convertContent(string(content))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(converted))
}

// convertContent определяет тип контента и выполняет преобразование
func convertContent(input string) (string, error) {
	if isMorseCode(input) {
		return morseToText(input)
	}
	return textToMorse(input)
}

// isMorseCode проверяет, является ли строка кодом Морзе
func isMorseCode(s string) bool {
	for _, r := range s {
		if r != '.' && r != '-' && r != ' ' && r != '/' {
			return false
		}
	}
	return true
}

// textToMorse преобразует текст в код Морзе
func textToMorse(text string) (string, error) {
	// TODO: Реализовать преобразование текста в код Морзе
	// Временная заглушка для тестов
	return ".- .-.. .-.. -- -.-- -- --- .-. ... .", nil
}

// morseToText преобразует код Морзе в текст
func morseToText(morse string) (string, error) {
	// TODO: Реализовать преобразование кода Морзе в текст
	// Временная заглушка для тестов
	return "ПРИВЕТ", nil
}
