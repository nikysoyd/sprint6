package server

import (
	"log"
	"net/http"
	"strings"
)

// Morse-код для русских букв (упрощенный пример)
var morseCode = map[rune]string{
	'А': ".-", 'Б': "-...", 'В': ".--", 'Г': "--.", 'Д': "-..",
	'Е': ".", 'Ж': "...-", 'З': "--..", 'И': "..", 'Й': ".---",
	'К': "-.-", 'Л': ".-..", 'М': "--", 'Н': "-.", 'О': "---",
	'П': ".--.", 'Р': ".-.", 'С': "...", 'Т': "-", 'У': "..-",
	'Ф': "..-.", 'Х': "....", 'Ц': "-.-.", 'Ч': "---.", 'Ш': "----",
	'Щ': "--.-", 'Ъ': "--.--", 'Ы': "-.--", 'Ь': "-..-", 'Э': "..-..",
	'Ю': "..--", 'Я': ".-.-",
}

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
		text := r.URL.Query().Get("text")
		var morse strings.Builder
		for _, char := range strings.ToUpper(text) {
			if code, ok := morseCode[char]; ok {
				morse.WriteString(code + " ")
			}
		}
		w.Write([]byte(morse.String()))
	})

	mux.HandleFunc("/morse-to-text", func(w http.ResponseWriter, r *http.Request) {
		// Реализуйте обратное преобразование (морзе → текст)
	})

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		// Обработка загрузки файла/текста
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
