package auth

import (
	"encoding/gob"
	"generic_inventory/conf"
	"generic_inventory/web"
	"net/http"
	"strings"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// User holds a users account information
type User struct {
	Username      string
	Authenticated bool
}

// Provider - An interface to allow interoperability with external authentication providers
type Authorization interface {
	CheckPass(name, password string) bool
}

func

// Store - holds all session data
var Store *sessions.CookieStore

func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	Store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	Store.Options = &sessions.Options{
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

	gob.Register(User{})
}

// GetUser returns a user from session s
// on error returns an empty user
func GetUser(s *sessions.Session) User {
	val := s.Values["user"]
	var user = User{}
	user, ok := val.(User)
	if !ok {
		return User{Authenticated: false}
	}
	return user
}

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

	// Where authentication could be done
	//if r.FormValue("code") != "code" {
	//	if r.FormValue("code") == "" {
	//		session.AddFlash("Must enter a code")
	//	}
	//	session.AddFlash("The code was incorrect")
	//	err = session.Save(r, w)
	//	if err != nil {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//	http.Redirect(w, r, "/forbidden", http.StatusFound)
	//	return
	//}
	var user User
	user.Username = r.FormValue("username")
	user.Authenticated, _ = provider.IsValid(user.Username, r.FormValue("code"))
	//user := &User{
	//	Username:      username,
	//	Authenticated: false,
	//}

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
