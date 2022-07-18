package users

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"

	ldap "github.com/go-ldap/ldap/v3"
)

func LDAPConn() (*ldap.Conn, error) {
	// TODO: env var/read from file
	// /etc/opt/oci-hpc/passwords/openldap/root.txt'
	password := "OglY2c8xgAfxAERVVQEM"
	tlsBastionURL := "bastion.cluster"
	tlsConfig := &tls.Config{ServerName: tlsBastionURL}
	l, err := ldap.DialURL("ldaps://localhost:636", ldap.DialWithTLSConfig(tlsConfig))
	l.Debug.Enable(true)
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
	userInfo := UserInfo{Username: login.Username}

	// userPassword - exists on the object for a given CN (a username) and OU (Organizational unit)
	pw, err := queryLDAPUserAttribute(PeopleDN, userInfo.Username, "userPassword")
	if err != nil {
		return userInfo, false
	}

	if pw == login.Password {
		return userInfo, true
	}

	return userInfo, false
}

func storeRefreshTokenLDAP(username string, refreshToken string) error {
	l, err := LDAPConn()
	if err != nil {
		fmt.Println("storeRefreshToken: " + err.Error())
		return err
	}
	defer l.Close()
	dn := fmt.Sprintf("CN=%s,ou=People,dc=local", username)
	mod := ldap.NewModifyRequest(dn, nil)
	mod.Replace(RefreshTokenLDAPKey, []string{refreshToken})
	if err := l.Modify(mod); err != nil {
		log.Println("storeRefreshToken: " + err.Error())
		return err
	}

	return nil
}

func queryLDAPUserAttribute(dn string, username string, attribute string) (string, error) {
	l, err := LDAPConn()
	if err != nil {
		log.Println("queryLDAPAttribute: " + err.Error())
		return "", err
	}
	defer l.Close()
	baseDN := dn
	// Filters must start and finish with ()
	filter := fmt.Sprintf("(CN=%s)", ldap.EscapeFilter(username))

	searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, filter, []string{attribute}, nil)

	result, err := l.Search(searchReq)
	if err != nil {
		log.Println("queryLDAPAttribute: " + err.Error())
		return "", err
	}
	if len(result.Entries) == 0 || len(result.Entries[0].Attributes) == 0 || len(result.Entries[0].Attributes[0].Values) == 0 {
		err = errors.New("queryLDAPAttribute: Could not find attribute")
		log.Println("queryLDAPAttribute: " + err.Error())
		return "", err
	}
	return result.Entries[0].Attributes[0].Values[0], err
}

func deleteDN(dn string) error {
	l, err := LDAPConn()
	if err != nil {
		log.Println("deleteDN: " + err.Error())
		return err
	}
	defer l.Close()
	delReq := ldap.NewDelRequest(dn, nil)
	if err := l.Del(delReq); err != nil {
		log.Println("deleteDN: " + err.Error())
		return err
	}
	return nil
}

// TODO:
// func createUser
// func updateUser
// func getUserInfo
