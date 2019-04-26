package api

import (
	"encoding/json"
	"fmt"
	"generic_inventory/auth"
	"generic_inventory/conf"
	"generic_inventory/dao"
	"generic_inventory/web"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// CrudHandle - Function to call other crud funtions
func CrudHandle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

	switch params["action"] {
	case "show":
		if params["sku"] == "all" {
			results, err := dao.GetEntries(conf.MyConfig)
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

// StaticHandle - Function to return static web components
func StaticHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	params := mux.Vars(r)
	switch params["path"] {
	case "css":
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, conf.MyConfig.StaticPath+params["path"]+"/"+params["file"]+"."+params["ext"])
	case "js":
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, conf.MyConfig.StaticPath+params["path"]+"/"+params["file"]+"."+params["ext"])
	case "img":
		w.Header().Set("Content-Type", "image/"+params["ext"])
		http.ServeFile(w, r, conf.MyConfig.StaticPath+params["path"]+"/"+params["file"]+"."+params["ext"])
	default:
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(404)
	}

}

// GetIndex - Return the main HTML page for the site
func GetIndex(w http.ResponseWriter, r *http.Request) {
	var tmpl = "login.html"
	p := &web.Page{Title: "Log in required", Body: []byte("This is a sample Page.")}

	session, err := auth.Store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := auth.GetUser(session)
	var permitted = "False"
	if user.Authenticated {
		permitted = "True"
	}
	fmt.Printf("%s %s %s \n", user.Username, user.Role, permitted)
	if user.Authenticated {
		tmpl = "index.html"
		p = &web.Page{Title: "Welcome To Generic Inventory", Body: []byte("This is a sample Page.")}
	}

	w.Header().Set("Content-Type", "text/html")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	web.RenderTemplate(w, tmpl, p, conf.MyConfig.TmplPath)
}

// RedirectToTLS - Handler for HTTp to HTTPS Redirection
func RedirectToTLS(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + req.URL.Path

	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target,
		// see @andreiavrammsd comment: often 307 > 301
		http.StatusTemporaryRedirect)
}
