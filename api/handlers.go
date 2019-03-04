package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// GetEntries - Return a json object containing people
func GetEntries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		panic(err)
	}
}

// GetEntry - Return a json object containing one person
func GetEntry(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, entry := range entries {
		if entry.SKU == params["sku"] {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(entry); err != nil {
				panic(err)
			}
			return
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
}

// CreateEntry - Create a json object containing one person
func CreateEntry(w http.ResponseWriter, r *http.Request) {
	var entry Entry
	params := mux.Vars(r)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &entry); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	entry.SKU = params["sku"]
	entries = append(entries, entry)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(entry); err != nil {
		panic(err)
	}
}

// DeleteEntry - Delete a person
func DeleteEntry(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, entry := range entries {
		if entry.SKU == params["sku"] {
			entries = append(entries[:index], entries[index+1:]...)
			break
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(entry); err != nil {
			panic(err)
		}
	}
}
