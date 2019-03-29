package main

import (
	"generic_inventory/api"
	"log"
	"net/http"
)

// main - our main function
func main() {
	router := api.NewRouter()
	log.Fatal(http.ListenAndServe(":8000", router))
}
