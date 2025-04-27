package handlers

import (
	//"html/template"
	"io"
	"strings"

	//"log"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/nikysoyd/sprint6/internal/service"
)

// ServeHome обрабатывает корневой эндпоинт и отдает index.html
func ServeHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	http.ServeFile(w, r, "sprint6/index.html")
}

// UploadHandler обрабатывает загрузку файла
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Ошибка при разборе формы", http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Не удалось получить файл из формы", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Не удалось прочитать файл", http.StatusInternalServerError)
		return
	}

	converted, err := service.Convert(string(data))
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка конвертации: %v", err), http.StatusInternalServerError)
		return
	}

	// Генерация имени файла с таймстемпом
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("converted_%s%s", strings.ReplaceAll(time.Now().UTC().Format(time.RFC3339), ":", "-"), ext)

	outFile, err := os.Create(filepath.Join(".", filename))

	if err != nil {
		http.Error(w, "Не удалось создать файл", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	if _, err := outFile.Write([]byte(converted)); err != nil {
		http.Error(w, "Не удалось записать результат", http.StatusInternalServerError)
		return
	}

	// Возвращаем пользователю результат
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(converted))
	w.WriteHeader(http.StatusOK)
}
