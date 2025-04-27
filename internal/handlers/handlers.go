package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nikysoyd/sprint6/internal/service"
)

func RootServeHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	http.ServeFile(w, r, "index.html")
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Form parsing error", http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Bad import from file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Bad file reading", http.StatusInternalServerError)
		return
	}

	converted, err := service.Convert(string(data))
	if err != nil {
		http.Error(w, fmt.Sprintf("Wrong convertation: %v", err), http.StatusInternalServerError)
		return
	}

	ext := filepath.Ext(header.Filename)
	newFileName := fmt.Sprintf("converted_%s%s", strings.ReplaceAll(time.Now().UTC().Format(time.RFC3339), ":", "-"), ext)

	outFile, err := os.Create(filepath.Join(".", newFileName))

	if err != nil {
		http.Error(w, "Cant create file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	if _, err := outFile.Write([]byte(converted)); err != nil {
		http.Error(w, "Cant write result", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(converted))
	w.WriteHeader(http.StatusOK)
}
