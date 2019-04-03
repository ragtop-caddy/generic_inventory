package web

import (
	"net/http"
	"text/template"
)

// RenderTemplate - Function to render standard templates
func RenderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles("C:/Users/dusti/go/src/generic_inventory/web/static/templates/index.html")
	t.Execute(w, p)
}
