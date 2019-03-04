package main

import "net/http"

// GetIndex - Return the main HTML page for the site
func GetIndex(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
