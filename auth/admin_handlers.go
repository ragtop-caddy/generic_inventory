package auth

import (
	"encoding/json"
	"generic_inventory/conf"
	"generic_inventory/web"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// AdminModel - An object to hold various administrative attributes
type AdminModel struct {
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

// AdminPanel - HTTP Handler to show the admin panel
func AdminPanel(w http.ResponseWriter, r *http.Request) {
	var tmpl = "admin.html"
	p := &web.Page{Title: "Welcome to the Admin Panel", Message: "You are an Admin"}
	session, err := Store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := GetUser(session)
	if user.Role != "admin" {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	w.Header().Set("Content-Type", "text/html")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	web.RenderTemplate(w, tmpl, p, conf.MyConfig.TmplPath)
}

// RetrieveUser - Return a JSON object containing the credentials for one user
func RetrieveUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	session, err := Store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := GetUser(session)
	if user.Role != "admin" {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	result, err := SessionAuth.retrieveUser(params["uid"])
	if err != nil {
		w.WriteHeader(404) // Not Found
	} else {
		if err := json.NewEncoder(w).Encode(result); err != nil {
			w.WriteHeader(500) // Internal error
		}
	}
}

// RetrieveUsers - Return a JSON object containing the credentials for one user
func RetrieveUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	result, err := SessionAuth.retrieveUsers()
	if err != nil {
		w.WriteHeader(404) // Not Found
	} else {
		if err := json.NewEncoder(w).Encode(result); err != nil {
			w.WriteHeader(500) // Internal error
		}
	}
}

// CreateUser - Creates a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var creds credentials
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}
	err = json.Unmarshal(body, &creds)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}
	result, err := SessionAuth.createUser(creds)
	if err != nil {
		log.Printf("ERROR: received %s while adding new user", err)
	}
	log.Printf("INFO: Got %s while creating new user %s", result, creds.UID)
}
