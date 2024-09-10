package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Handle routes for serving files and uploading files
	http.HandleFunc("/static", serveStaticFile)
	http.HandleFunc("/upload", uploadFile)

	// Start the server
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// serveStaticFile serves a static file from the 'static' directory based on the 'file' query parameter.
func serveStaticFile(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	staticDir := "./static"
	filePath := filepath.Join(staticDir, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}

// uploadFile handles file uploads and saves them to the 'static' directory.
func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form with a max upload size of 10MB
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	// Retrieve the file from form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Could not retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create the directory if it doesn't exist
	staticDir := "./static"
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		err = os.Mkdir(staticDir, os.ModePerm)
		if err != nil {
			http.Error(w, "Could not create directory", http.StatusInternalServerError)
			return
		}
	}

	// Create the file on the server
	filePath := filepath.Join(staticDir, handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Could not create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file's content to the server file
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("File uploaded successfully"))
}
