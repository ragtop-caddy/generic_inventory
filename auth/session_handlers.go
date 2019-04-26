package auth

import (
	"fmt"
	"generic_inventory/conf"
	"generic_inventory/web"
	"net/http"
	"strings"
)

// HTTP Handlers

// ValidateSession - Function to track session state for client connections
func ValidateSession(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var path string
		if r.URL.Path != "/" {
			path = strings.Split(r.URL.Path, "/")[1]
		} else {
			path = r.URL.Path
		}
		unprotected := map[string]bool{
			"forbidden": true,
			"login":     true,
			"logout":    true,
			"css":       true,
			"img":       true,
		}
		session, err := Store.Get(r, "cookie-name")
		if session.IsNew {
			if err != nil {
				fmt.Printf("INFO: Found bogus cookie\n")
				session.Options.MaxAge = -1
				err = session.Save(r, w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				//} else {
				//	fmt.Printf("INFO: New user session, saving cookie\n")
				//	//err = session.Save(r, w)
				//	if err != nil {
				//		http.Error(w, err.Error(), http.StatusInternalServerError)
				//	}
			}
		}
		user := GetUser(session)
		if !unprotected[path] {
			if !user.Authenticated && r.URL.Path != "/" {
				http.Redirect(w, r, "/", http.StatusFound)
			}
			if user.Authenticated {
				err = session.Save(r, w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
		inner.ServeHTTP(w, r)
	})
}

// Forbidden - http Handler to render the login page
func Forbidden(w http.ResponseWriter, r *http.Request) {
	var tmpl = "login.html"
	p := &web.Page{Title: "Log in required", Body: []byte("This is a sample Page.")}
	w.Header().Set("Content-Type", "text/html")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	web.RenderTemplate(w, tmpl, p, conf.MyConfig.TmplPath)
}

// Login - authenticates the user
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	session, err := Store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := SessionAuth.Authenticate(r.FormValue("username"), r.FormValue("code"))
	session.Values["user"] = user
	fmt.Printf("INFO: Login got %s, saving cookie \n", session.Values["user"])
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// Logout - revokes authentication for a user
func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	session, err := Store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = User{}
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
