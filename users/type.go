package users

import "time"

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

type RefreshToken struct {
	RefreshTokenString string
	Expiration         int64
}

const (
	AccessTokenKey  = "access_token"
	RefreshTokenKey = "refresh_token"
	//TODO: This should be moved to a custom key
	// However, that requires a specific change to LDAP config
	RefreshTokenLDAPKey          = "description"
	AccessTokenExpirationWindow  = 1 * time.Minute
	RefreshTokenExpirationWindow = 7 * 24 * time.Hour
	// check ldap information with `sudo slapcat` command
	// ou=People,DC=local - default location for users
	BaseDN = "ou=People,DC=local"
)
