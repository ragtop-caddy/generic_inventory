package auth

import (
	"generic_inventory/conf"
	"generic_inventory/web"
	"net/http"
	"text/template"
)

// AdminModel - An object to hold various administrative attributes
type AdminModel struct {
	Title      string
	Body       []byte
	UID        string `json:"uid,omitempty" bson:"uid,omitempty"`
	Fname      string `json:"fname,omitempty" bson:"fname,omitempty"`
	Lname      string `json:"lname,omitempty" bson:"lname,omitempty"`
	Email      string `json:"email,omitempty" bson:"email,omitempty"`
	Role       string `json:"role,omitempty" bson:"role,omitempty"`
	State      string `json:"state,omitempty" bson:"state,omitempty"`
	Inactive   int    `json:"inactive,omitempty" bson:"inactive,omitempty"`
	Expiration int    `json:"expiration,omitempty" bson:"expiration,omitempty"`
	Last       string `json:"last,omitempty" bson:"last,omitempty"`
}

// RenderTemplate - Function to render standard templates
func renderTemplate(w http.ResponseWriter, tmpl string, p *AdminModel, tmplPath string) {
	t, _ := template.ParseFiles(tmplPath + tmpl)
	t.Execute(w, p)
}

// AdminPanel - HTTP Handler to show the admin panel
func AdminPanel(w http.ResponseWriter, r *http.Request) {
	var tmpl = "admin.html"
	p := &web.Page{Title: "Welcome to the Admin Panel", Message: "You are not an Admin"}
	session, err := Store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := GetUser(session)
	if user.Role == "admin" {
		p.Message = "You are an Admin"
	}
	w.Header().Set("Content-Type", "text/html")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	web.RenderTemplate(w, tmpl, p, conf.MyConfig.TmplPath)
}

// RetrieveUser - HTTP Handler for retrieving user information
func RetrieveUser(w http.ResponseWriter, r *http.Request) {
	var tmpl = "login.html"
	p := &web.Page{Title: "Log in required", Body: []byte("This is a sample Page.")}
	w.Header().Set("Content-Type", "text/html")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	web.RenderTemplate(w, tmpl, p, conf.MyConfig.TmplPath)
}
