package server // Измените package на main, так как у вас есть функция main()

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/nikysoyd/sprint6/internal/service" // Исправленный импорт
)

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/upload", uploadHandler)

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Улучшенный путь к файлу
	html, err := os.ReadFile("./static/index.html") // Поместите index.html в папку static
	if err != nil {
		http.Error(w, "Не удалось прочитать index.html: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(html)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Ошибка при разборе формы: "+err.Error(), http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Ошибка при получении файла: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Ошибка при чтении файла: "+err.Error(), http.StatusInternalServerError)
		return
	}

	converted, err := service.AutoDetectAndConvert(string(fileContent))
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при конвертации: %v", err), http.StatusInternalServerError)
		return
	}

	// Создаем папку results, если её нет
	if err := os.MkdirAll("results", 0755); err != nil {
		http.Error(w, "Ошибка при создании папки: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fileName := filepath.Join("results", fmt.Sprintf("result_%s%s",
		time.Now().UTC().Format("20060102150405"),
		filepath.Ext(header.Filename)))

	outputFile, err := os.Create(fileName)
	if err != nil {
		http.Error(w, "Ошибка при создании файла результата: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	if _, err = outputFile.WriteString(converted); err != nil {
		http.Error(w, "Ошибка при записи результата: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(converted))
}
