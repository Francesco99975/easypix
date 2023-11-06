package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	controller "github.com/Francesco99975/easypix/internal"
)


func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Fprintf(w, "Nothing to see here...")
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/upload", controller.Upload)
	http.HandleFunc("/delete/", controller.Delete)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("uploads"))))

	err := os.Mkdir("uploads", os.ModePerm)
	if err != nil {
		fmt.Printf("Could not create uploads directory. Error: %s.", err.Error())
	}
	
	fmt.Println("Server is running on :8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}