package auth

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// User holds a users account information
type User struct {
	Username      string
	Authenticated bool
	Role          string
}

// Authenticator - An interface to allow interoperability with external authentication providers
type Authenticator interface {
	Authenticate(r *http.Request) *User
}

// Manager - An interface to allow user management tasks via additional providers
type Manager interface {
	ShowUser(r *http.Request)
	CreateUser(r *http.Request)
	DeleteUser(r *http.Request)
	UpdateUser(r *http.Request)
}

// SessionAuth - Object to interact with the internal auth module
var SessionAuth InternalAuth

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
