package auth

import (
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
		if err != nil && session.IsNew {
			session.Options.MaxAge = -1
			err = session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		if !unprotected[path] {
			user := GetUser(session)
			if !user.Authenticated {
				http.Redirect(w, r, "/forbidden", http.StatusFound)
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

	user := sessionAuth.CheckPass(r)
	session.Values["user"] = user

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
