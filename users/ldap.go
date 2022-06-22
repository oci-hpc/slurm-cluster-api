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

// Validate checks whether a user/password combination is valid
func Validate(login LoginInfo) bool {
	// TODO: return true for now; plug in ldap later
	return true
}
