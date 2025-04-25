package handlers

import (
	"log"
	"net/http"
	"os"
	"time"
)

// Server структура для нашего HTTP сервера
type Server struct {
	logger     *log.Logger
	httpServer *http.Server
}

// NewServer создает новый экземпляр сервера
func NewServer(logger *log.Logger) *Server {
	// Создаем роутер
	router := createRouter(logger)

	// Конфигурируем HTTP сервер
	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		logger:     logger,
		httpServer: httpServer,
	}
}

// createRouter создает и настраивает HTTP роутер
func createRouter(logger *log.Logger) *http.ServeMux {
	router := http.NewServeMux()

	// Регистрируем обработчики
	router.HandleFunc("/", rootHandler(logger))
	router.HandleFunc("/upload", uploadHandler(logger))

	return router
}

// rootHandler обработчик для корневого пути
func rootHandler(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		logger.Println("Handling root request")
		w.Write([]byte("Hello, World!"))
	}
}

// uploadHandler обработчик для загрузки файлов
func uploadHandler(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		logger.Println("Handling upload request")
		// Здесь должна быть логика обработки загрузки файла
		w.Write([]byte("File uploaded successfully"))
	}
}

// Start запускает сервер
func (s *Server) Start() error {
	s.logger.Printf("Starting server on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func main() {
	// Инициализируем логгер
	logger := log.New(os.Stdout, "HTTP_SERVER: ", log.LstdFlags|log.Lshortfile)

	// Создаем и запускаем сервер
	server := NewServer(logger)
	if err := server.Start(); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
