package main

import (
	"log"
	"net/http"
)

// main - our main function
func main() {
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8000", router))
}
