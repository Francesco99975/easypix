package controller

import (
	"io"
	"net/http"
	"os"
	"path"
)


func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error while retriving file", http.StatusInternalServerError)
		}
		defer file.Close()

		uploadedFile, err := os.Create(path.Join("uploads", handler.Filename))
		if err != nil {
			http.Error(w, "Error while writing file on the server", http.StatusInternalServerError)
		}
		defer uploadedFile.Close()

		_, err = io.Copy(uploadedFile, file)
		if err != nil {
			http.Error(w, "Error while saving the file", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Invalid Method. Onlu POST Allowed.", http.StatusMethodNotAllowed)
	}
}