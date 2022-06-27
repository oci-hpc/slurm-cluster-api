package users

import (
	"crypto/tls"
	"fmt"

	ldap "github.com/go-ldap/ldap/v3"
)

func LDAPConn() (*ldap.Conn, error) {
	// TODO: env var/read from file
	// /etc/opt/oci-hpc/passwords/openldap/root.txt'
	//uri := "ldaps://localhost:686"
	//tlsBastionURL := "bastion.cluster"
	password := "OglY2c8xgAfxAERVVQEM"
	tlsBastionURL := "bastion.cluster"
	tlsConfig := &tls.Config{ServerName: tlsBastionURL}
	l, err := ldap.DialURL("ldaps://localhost:636", ldap.DialWithTLSConfig(tlsConfig))

	if err != nil {
		return l, err
	}

	err = l.Bind("cn=manager,dc=local", password)
	if err != nil {
		return l, err
	}
	return l, err
}

// validateLDAPLogin checks whether a user/password combination is valid
func validateLDAPLogin(login LoginInfo) (UserInfo, bool) {
	// TODO: always return true for now; should validate via LDAP
	// TODO: pull username and info from LDAP
	userInfo := UserInfo{Username: login.Username}
	l, err := LDAPConn()
	if err != nil {
		fmt.Println(err.Error())
		return userInfo, false
	}
	defer l.Close()
	// check ldap information with `sudo slapcat` command
	// ou=People,DC=local - default location for users
	baseDN := "ou=People,DC=local"
	// Filters must start and finish with ()
	filter := fmt.Sprintf("(CN=%s)", ldap.EscapeFilter(userInfo.Username))

	// userPassword - exists on the object for a given CN (a username) and OU (Organizational unit)
	searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, filter, []string{"userPassword"}, nil)

	result, err := l.Search(searchReq)
	if err != nil {
		fmt.Println(err.Error())
		return userInfo, false
	}
	if len(result.Entries) == 0 || len(result.Entries[0].Attributes) == 0 || len(result.Entries[0].Attributes[0].Values) == 0 {
		return userInfo, false
	}
	pw := result.Entries[0].Attributes[0].Values[0]

	if pw == login.Password {
		return userInfo, true
	}

	return userInfo, false
}

// TODO:
// func createUser
// func updateUser
// func getUserInfo
