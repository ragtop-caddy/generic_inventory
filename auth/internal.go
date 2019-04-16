package auth

import "net/http"

type internalUser struct {
	name     string
	password string
}

// InternalUsers - struct containing a list of internal users
type InternalUsers []internalUser

// NewInternalAuth - Create a new internal auth mechanism
func NewInternalAuth() InternalUsers {
	var u = InternalUsers{
		internalUser{name: "joe", password: "code"},
		internalUser{name: "bill", password: "foo"},
	}
	return u
}

// CheckPass - Check password validity
func (iu InternalUsers) CheckPass(r *http.Request) User {
	var u User
	u.Username = r.FormValue("username")
	u.Authenticated = false
	for _, user := range iu {
		if user.name == u.Username {
			if r.FormValue("code") == user.password {
				u.Authenticated = true
			}
		}
	}
	return u
}
