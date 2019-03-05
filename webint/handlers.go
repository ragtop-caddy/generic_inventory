package main

import "net/http"

var p Page

// GetIndex - Return the main HTML page for the site
func GetIndex(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Hello", Body: []byte("This is a sample Page.")}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	renderTemplate(w, "view", p)
}
