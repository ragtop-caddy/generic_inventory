package main

import (
	"generic_inventory/api"
	"log"
	"net/http"
)

// main - our main function
func main() {
	var dao api.InventoryDAO
	dao.URI = "mongodb://localhost:27017"
	dao.ConfigDB()
	router := api.NewRouter()
	log.Fatal(http.ListenAndServe(":8000", router))
}
