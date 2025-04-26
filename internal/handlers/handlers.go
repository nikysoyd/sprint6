package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"your_module_name/service" // Замените на актуальный путь к пакету service
)

// IndexHandler обрабатывает запрос к корневому эндпоинту /
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// UploadHandler обрабатывает загрузку файла и конвертацию
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсинг формы с лимитом 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получение файла из формы
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Чтение содержимого файла
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Конвертация данных с помощью функции из пакета service
	result, err := service.Convert(string(data))
	if err != nil {
		http.Error(w, "Conversion failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Генерация имени выходного файла
	ext := filepath.Ext(handler.Filename)
	outputFilename := fmt.Sprintf("output_%s%s", time.Now().UTC().Format("20060102150405"), ext)
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		http.Error(w, "Failed to create output file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	// Запись результата в файл
	_, err = outputFile.WriteString(result)
	if err != nil {
		http.Error(w, "Failed to write to output file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка результата клиенту
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Conversion successful. Result: %s\nSaved to: %s", result, outputFilename)
}
