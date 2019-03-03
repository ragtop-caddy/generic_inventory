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
	Amount   int    `json:"amount,omitempty"`
}

// Header - Information standard to all entries
type Header struct {
	Type        string        `json:"type,omitempty"`
	Description string        `json:"description,omitempty"`
	Stock       int           `json:"stock,omitempty"`
	History     []Transaction `json:"history,omitempty"`
}

// Detail - Defines entry details
type Detail struct {
	Gender string `json:"gender,omitempty"`
	Color  string `json:"color,omitempty"`
	Size   string `json:"size,omitempty"`
	Style  string `json:"style,omitempty"`
	Fit    string `json:"fit,omitempty"`
}

// Entry - Defines various types of inventory
type Entry struct {
	SKU     string  `json:"sku,omitempty"`
	Header  *Header `json:"header,omitempty"`
	Details *Detail `json:"details,omitempty"`
}

var entries []Entry

// GetEntries - Return a json object containing people
func GetEntries(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(entries)
}

// GetEntry - Return a json object containing one person
func GetEntry(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range entries {
		if item.SKU == params["sku"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

// CreateEntry - Create a json object containing one person
func CreateEntry(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var entry Entry
	_ = json.NewDecoder(r.Body).Decode(&entry)
	entry.SKU = params["sku"]
	entries = append(entries, entry)
	json.NewEncoder(w).Encode(entry)
}

// DeleteEntry - Delete a person
func DeleteEntry(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range entries {
		if item.SKU == params["sku"] {
			entries = append(entries[:index], entries[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(item)
	}
}

// main - our main function
func main() {
	router := mux.NewRouter()

	entries = append(entries, Entry{
		SKU: "1",
		Header: &Header{
			Type:        "pant",
			Description: "boys small pant",
			Stock:       5,
			History: []Transaction{
				Transaction{
					ISODate:  "1234",
					Campus:   "duffy",
					Students: 0,
					Action:   "create",
				},
			},
		},
		Details: &Detail{
			Gender: "m",
			Color:  "blk",
			Size:   "small",
			Style:  "slacks",
			Fit:    "loose",
		},
	},
	)

	router.HandleFunc("/inventory", GetEntries).Methods("GET")
	router.HandleFunc("/inventory/{sku}", GetEntry).Methods("GET")
	router.HandleFunc("/inventory/{sku}", CreateEntry).Methods("POST")
	router.HandleFunc("/inventory/{sku}", DeleteEntry).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}
