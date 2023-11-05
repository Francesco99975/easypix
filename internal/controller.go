package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type UploadResponse struct {
	URL string `json:"url"`
}


func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error while retriving file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		ext := filepath.Ext(handler.Filename)
		allowedExts := map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".webp": true}
		if !allowedExts[ext] {
			http.Error(w, "Unsupported file type. Only PNG, JPG, JPEG, and WEBP are allowed.", http.StatusBadRequest)
			return
		}

		uniqueFilename := genid(ext)
		uploadedFile, err := os.Create(path.Join("uploads", uniqueFilename))
		if err != nil {
			http.Error(w, "Error while writing file on the server", http.StatusInternalServerError)
			return
		}
		defer uploadedFile.Close()

		_, err = io.Copy(uploadedFile, file)
		if err != nil {
			http.Error(w, "Error while saving the file", http.StatusInternalServerError)
			return
		}

		url := fmt.Sprintf("/images/%s", uniqueFilename)
		response := UploadResponse{URL: url}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, `{"error": "Error generating JSON response"}`, http.StatusInternalServerError)
			return
		}

		w.Write(jsonResponse)
	} else {
		http.Error(w, "Invalid Method. Onlu POST Allowed.", http.StatusMethodNotAllowed)
	}
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File parameter is missing", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("uploads", fileName)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set the Content-Disposition header to suggest a filename for the downloaded file
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	// Determine the file's content type based on the file extension
	contentType := "application/octet-stream"
	switch strings.ToLower(filepath.Ext(fileName)) {
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".webp":
		contentType = "image/webp"
	}

	w.Header().Set("Content-Type", contentType)

	// Serve the file
	io.Copy(w, file)
}
