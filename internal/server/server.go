package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"your_project/service" // импортируем ваш пакет с функцией автоопределения
)

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/upload", uploadHandler)

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// rootHandler обрабатывает корневой эндпоинт /
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Читаем и возвращаем содержимое index.html
	html, err := os.ReadFile("index.html")
	if err != nil {
		http.Error(w, "Не удалось прочитать index.html", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(html)
}

// uploadHandler обрабатывает эндпоинт /upload
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Парсим форму, включая загруженный файл
	err := r.ParseMultipartForm(10 << 20) // Максимум 10MB файл
	if err != nil {
		http.Error(w, "Ошибка при разборе формы", http.StatusInternalServerError)
		return
	}

	// Получаем файл из формы
	file, header, err := r.FormFile("file") // "file" - имя поля в форме
	if err != nil {
		http.Error(w, "Ошибка при получении файла", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Читаем содержимое файла
	fileContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Ошибка при чтении файла", http.StatusInternalServerError)
		return
	}

	// Конвертируем содержимое с помощью вашей функции из пакета service
	converted, err := service.AutoDetectAndConvert(string(fileContent))
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при конвертации: %v", err), http.StatusInternalServerError)
		return
	}

	// Создаем файл для сохранения результата
	fileName := fmt.Sprintf("result_%s%s",
		time.Now().UTC().Format("20060102150405"),
		filepath.Ext(header.Filename))

	outputFile, err := os.Create(fileName)
	if err != nil {
		http.Error(w, "Ошибка при создании файла результата", http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	// Записываем результат в файл
	_, err = outputFile.WriteString(converted)
	if err != nil {
		http.Error(w, "Ошибка при записи результата", http.StatusInternalServerError)
		return
	}

	// Возвращаем результат клиенту
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(converted))
}
