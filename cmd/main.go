package main

import (
	"log"
	"net/http"
)


func main() {
	http.Handle("/upload", controller.upload)
	

	log.Fatal(http.ListenAndServe(":8888", nil))
}