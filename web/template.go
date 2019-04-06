package web

import (
	"net/http"
	"text/template"
)

// RenderTemplate - Function to render standard templates
func RenderTemplate(w http.ResponseWriter, tmpl string, p *Page, tmplPath string) {
	t, _ := template.ParseFiles(tmplPath + tmpl)
	t.Execute(w, p)
}
