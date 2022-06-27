package users

import (
	"crypto/tls"
	"log"

	ldap "github.com/go-ldap/ldap/v3"
)

func LDAPConn() {
	// TODO: env var/read from file
	// /etc/opt/oci-hpc/passwords/openldap/root.txt'
	uri := "bastion:636"
	tlsBastionURL := "bastion.cluster"
	password := "VBNF4144Kl8C0qGD8xHa"

	tlsConfig := &tls.Config{ServerName: tlsBastionURL}
	l, err := ldap.DialTLS("tcp", uri, tlsConfig)

	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	err = l.Bind("cn=manager,dc=local", password)
	if err != nil {
		log.Fatal(err)
	}
}

// validateLDAPLogin checks whether a user/password combination is valid
func validateLDAPLogin(login LoginInfo) (UserInfo, bool) {
	// TODO: always return true for now; should validate via LDAP
	// TODO: pull username and info from LDAP
	userInfo := UserInfo{Username: login.Username, Email: "test@test.com"}
	return userInfo, true
}

// TODO:
// func createUser
// func updateUser
// func getUserInfo
