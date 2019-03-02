package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Transaction - Standard transaction structure
type Transaction struct {
	ISODate  string `json:"isodate,omitempty"`
	Campus   string `json:"campus,omitempty"`
	Students int    `json:"students,omitempty"`
	Action   string `json:"action,omitempty"`
}

// EntryHeader - Information standard to all objects
type EntryHeader struct {
	SKU         string        `json:"sku,omitempty"`
	Type        string        `json:"type,omitempty"`
	Description string        `json:"description,omitempty"`
	Stock       int           `json:"stock,omitempty"`
	History     []Transaction `json:"history,omitempty"`
}

// GarmentDetail - Defines Apparel specific details
type GarmentDetail struct {
	Gender string `json:"gender,omitempty"`
	Color  string `json:"color,omitempty"`
	Size   string `json:"size,omitempty"`
	Style  string `json:"style,omitempty"`
	Fit    string `json:"fit,omitempty"`
}

// Garment - Defines various types of apparel
type Garment struct {
	Header  *EntryHeader   `json:"header,omitempty"`
	Details *GarmentDetail `json:"details,omitempty"`
}

var apparel []Garment

// GetApparel - Return a json object containing people
func GetApparel(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(apparel)
}

// GetGarment - Return a json object containing one person
func GetGarment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range apparel {
		if item.Header.SKU == params["sku"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

// CreateGarment - Create a json object containing one person
func CreateGarment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var garment Garment
	_ = json.NewDecoder(r.Body).Decode(&garment)
	garment.Header.SKU = params["header"]
	apparel = append(apparel, garment)
	json.NewEncoder(w).Encode(apparel)
}

// DeletePerson - Delete a person
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}

// our main function
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}
