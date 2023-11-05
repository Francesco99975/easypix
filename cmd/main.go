package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	controller "github.com/Francesco99975/easypix/internal"
)


func main() {
	
	http.HandleFunc("/upload", controller.Upload)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("uploads"))))

	os.Mkdir("uploads", os.ModePerm)
	
	fmt.Println("Server is running on :8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}