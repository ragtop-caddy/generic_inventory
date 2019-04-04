package main

import (
	"generic_inventory/api"
	"generic_inventory/dao"
	"log"
	"net/http"
	"os"
)

// main - our main function
func main() {
	var dao dao.InventoryDAO
	uri, ok := os.LookupEnv("MONGODB_URI")
	if ok {
		dao.URI = uri
	} else {
		dao.URI = "mongodb://localhost:27017"
	}
	dao.ConfigDB()
	router := api.NewRouter()
	log.Fatal(http.ListenAndServe(":8000", router))
}
