package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	convert "github.com/Francesco99975/easypix/pkg"
)

type UploadResponse struct {
	URL string `json:"url"`
}


func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "File too large", http.StatusBadRequest)
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

		var uniqueFilename string
		var uploadedFile *os.File
		if ext != ".webp" {
			uniqueFilename = genid(".webp")
			uploadedFile, err = os.Create(path.Join("uploads", uniqueFilename))
			if err != nil {
				http.Error(w, "Error while writing file on the server", http.StatusInternalServerError)
				return
			}

			err := convert.ToWebP(&file, uploadedFile)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
		} else {
			uniqueFilename = genid(ext)
			uploadedFile, err = os.Create(path.Join("uploads", uniqueFilename))
			if err != nil {
				http.Error(w, "Error while writing file on the server", http.StatusInternalServerError)
				return
			}

			_, err = io.Copy(uploadedFile, file)
			if err != nil {
				http.Error(w, "Error while saving the file", http.StatusInternalServerError)
				return
			}
		}
		defer uploadedFile.Close()

		

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