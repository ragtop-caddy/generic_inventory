package auth

import (
	"encoding/json"
	"generic_inventory/conf"
	"generic_inventory/dao"
	"generic_inventory/web"
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

// AdminCrudHandle - Function to call other crud funtions
func AdminCrudHandle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	switch params["action"] {
	case "show":
		if params["uid"] == "all" {
			creds, err := retrieveCreds()
			if err != nil {
				w.WriteHeader(404) // Not Found
			}
			if err := json.NewEncoder(w).Encode(results); err != nil {
				w.WriteHeader(500) // Internal error
			}
		} else {
			result, err := dao.GetEntry(params["sku"], conf.MyConfig)
			if err != nil {
				w.WriteHeader(404) // Not Found
			} else {
				if err := json.NewEncoder(w).Encode(result); err != nil {
					w.WriteHeader(500) // Internal error
				}
			}
		}
	case "add":
		id, err := dao.CreateEntry(r.Body, params["sku"], conf.MyConfig)
		if err != nil {
			w.WriteHeader(500) // Internal error
			if err := json.NewEncoder(w).Encode(err); err != nil {
				w.WriteHeader(500) // Internal error
			}
		} else {
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(id); err != nil {
				w.WriteHeader(500) // Internal error
			}
		}
	case "remove":
		count, err := dao.DeleteEntry(params["sku"], conf.MyConfig)
		if err != nil {
			w.WriteHeader(500) // Internal error
			if err := json.NewEncoder(w).Encode(err); err != nil {
				w.WriteHeader(500) // Internal error
			}
		} else {
			if count == 0 {
				w.WriteHeader(404) // Not Found
			}
			if err := json.NewEncoder(w).Encode(count); err != nil {
				w.WriteHeader(500) // Internal error
			}
		}
	default:
		if err := json.NewEncoder(w).Encode(r.URL.RawQuery); err != nil {
			w.WriteHeader(500) // Internal error
		}
	}
}
