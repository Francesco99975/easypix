package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

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

		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error while retriving file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil && err != io.EOF {
			http.Error(w, "Failed to read the file", http.StatusInternalServerError)
			return
		}

		_, err = file.Seek(0, 0)
		if err != nil {
			http.Error(w, "Failed to seek to start of the file", http.StatusInternalServerError)
			return
		}

		if !isFormatAllowed(buffer) {
			http.Error(w, "Unsupported file type. Only PNG, JPG, JPEG, and WEBP are allowed.", http.StatusBadRequest)
			return
		}

		uniqueFilename := genid(".webp")
		uploadedFile, err := os.Create(path.Join("uploads", uniqueFilename))
		if err != nil {
			http.Error(w, "Error while writing file on the server", http.StatusInternalServerError)
			return
		}

		if !magicMatches(buffer, []byte{'R', 'I', 'F', 'F'}) {
			err := convert.ToWebP(&file, uploadedFile)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
		} else {
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
		http.Error(w, "Invalid Method. Only POST Allowed.", http.StatusMethodNotAllowed)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		url := r.URL.Path
		segments := strings.Split(url, "/")

		if len(segments) >= 3 {
			id := segments[2]

			err := os.Remove(path.Join("uploads", fmt.Sprintf("%s.webp", id)))

			if err != nil {
				http.Error(w, "Image file not found...", http.StatusNotFound)
			}
		} else {
			http.Error(w, "Paramenter not provided...", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Invalid Method. Only DELETE Allowed.", http.StatusMethodNotAllowed)
	}
}

func isFormatAllowed(buffer []byte) bool  {
	jpegMagic := []byte{0xFF, 0xD8}
    pngMagic := []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}
    webpMagic := []byte{'R', 'I', 'F', 'F'}

	return magicMatches(buffer, jpegMagic) || magicMatches(buffer, pngMagic) || magicMatches(buffer, webpMagic)
}

func magicMatches(buffer, magic []byte) bool {
	return len(buffer) >= len(magic) && bytesEqual(buffer[:len(magic)], magic)
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	} else {
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}

		return true
	}
}