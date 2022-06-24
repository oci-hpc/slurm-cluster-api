package users

// LoginInfo contains user-provided information to authenticate a user
type LoginInfo struct {
	Username string
	Password string
}

// UserInfo contains user info augmented by LDAP
type UserInfo struct {
	Username string
	Email    string
	//TODO: pull in info from LDAP
}
