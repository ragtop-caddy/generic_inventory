package main

import (
	"html/template"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles("./templates/index.html")
	t.Execute(w, p)
}
