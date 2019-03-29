package web

import (
	"html/template"
	"net/http"
)

// RenderTemplate - Function to render standard templates
func RenderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles("./templates/index.html")
	t.Execute(w, p)
}
