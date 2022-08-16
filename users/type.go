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

type RoleRequest struct {
	Name string
	Role string
}

type ClaimRequest struct {
	Name  string
	Value int
	Role  string
}

type RefreshToken struct {
	RefreshTokenString string
	Expiration         int64
}

type RBACClaim struct {
	Name  string
	Value int
	DN    string
}

type RBACRole struct {
	Name   string
	Claims []RBACClaim
	Users  []UserInfo
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
	PeopleDN      = "ou=People,DC=local"
	ClaimDN       = "cn=cluster,dc=local"
	ClaimNameKey  = "x-cluster-claim-name"
	ClaimValueKey = "x-cluster-claim-value"
)
