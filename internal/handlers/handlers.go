package handlers

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/nikysoyd/sprint6/internal/service"
)

const (
	formFileName = "file"
	maxFileSize  = 1 << 20 // 1MB
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form
	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile(formFileName)
	if err != nil {
		log.Printf("Error getting file from form: %v", err)
		http.Error(w, "Error getting file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file content
	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Convert the content
	content := string(data)
	result, err := service.AutoDetectAndConvert(content)
	if err != nil {
		log.Printf("Error converting content: %v", err)
		http.Error(w, "Error converting content", http.StatusInternalServerError)
		return
	}

	// Create a local file with the result
	fileName := time.Now().UTC().Format("2006-01-02T15-04-05.999999999") + filepath.Ext(header.Filename)
	if err := os.WriteFile(fileName, []byte(result), 0644); err != nil {
		log.Printf("Error writing result file: %v", err)
		http.Error(w, "Error saving result", http.StatusInternalServerError)
		return
	}

	// Return the result to the client
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}
