package auth

// InternalCredentials - Struct to hold user credentials
type InternalCredentials struct {
	username string
	password string
}

func (ic InternalCredentials) CheckPass() bool {
	if ic.password == "code" {
		return true
	}
	return false
}
