package users

import (
	"errors"
	"fmt"
	"log"

	ldap "github.com/go-ldap/ldap/v3"
)

func CheckRoleExists(role string) bool {
	res, err := QueryRBACRole(role)
	if err != nil {
		return false
	}
	if len(res) == 0 {
		return false
	}
	return true
}

func QueryRBACRole(role string) ([]*ldap.Entry, error) {
	l, err := LDAPConn()
	if err != nil {
		log.Println("QueryRBACRole: " + err.Error())
		return nil, err
	}
	defer l.Close()
	// Filters must start and finish with ()
	filter := fmt.Sprintf("(cn=%s)", ldap.EscapeFilter(role))

	searchReq := ldap.NewSearchRequest(PeopleDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, filter, nil, nil)

	result, err := l.Search(searchReq)
	if err != nil {
		log.Println("QueryRBACRole: " + err.Error())
		return nil, err
	}
	if len(result.Entries) == 0 {
		err = errors.New("Could not find role")
		log.Println("QueryRBACRole: " + err.Error())
		return nil, err
	}
	return result.Entries, err
}

func QueryAllRBACClaims() ([]*ldap.Entry, error) {
	l, err := LDAPConn()
	if err != nil {
		log.Println("QueryRBACClaim: " + err.Error())
		return nil, err
	}
	defer l.Close()
	// Filters must start and finish with ()
	filter := fmt.Sprintf("(objectClass=x-cluster-claim)")

	searchReq := ldap.NewSearchRequest(ClaimDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, filter, nil, nil)

	result, err := l.Search(searchReq)
	if err != nil {
		log.Println("QueryRBACClaim: " + err.Error())
		return nil, err
	}
	if len(result.Entries) == 0 {
		err = errors.New("Could not find role")
		log.Println("QueryRBACClaim: " + err.Error())
		return nil, err
	}
	return result.Entries, err
}

func QueryRBACClaim(name string, value int) ([]*ldap.Entry, error) {
	l, err := LDAPConn()
	if err != nil {
		log.Println("QueryRBACClaim: " + err.Error())
		return nil, err
	}
	defer l.Close()
	// Filters must start and finish with ()
	filter := fmt.Sprintf("(cn=claim-%s-%d)", name, value)

	searchReq := ldap.NewSearchRequest(ClaimDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, filter, nil, nil)

	result, err := l.Search(searchReq)
	if err != nil {
		log.Println("QueryRBACClaim: " + err.Error())
		return nil, err
	}
	if len(result.Entries) == 0 {
		err = errors.New("Could not find role")
		log.Println("QueryRBACClaim: " + err.Error())
		return nil, err
	}
	return result.Entries, err
}

func QueryRBACRoleClaim(role string) ([]*ldap.Entry, error) {
	l, err := LDAPConn()
	if err != nil {
		log.Println("QueryRBACRoleClaim: " + err.Error())
		return nil, err
	}
	defer l.Close()
	// Filters must start and finish with ()
	filter := "(cn=cluster-claims)"
	//filter := fmt.Sprintf("(cn=%s)", ldap.EscapeFilter(role))
	dn := fmt.Sprintf("cn=%s,%s", role, PeopleDN)
	searchReq := ldap.NewSearchRequest(dn, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, filter, nil, nil)

	result, err := l.Search(searchReq)

	if err != nil {
		log.Println("QueryRBACRoleClaim: " + err.Error())
		return nil, err
	}

	return result.Entries, err
}

func QueryRBACRoles() ([]*ldap.Entry, error) {
	l, err := LDAPConn()
	if err != nil {
		log.Println("QueryRBACRoles: " + err.Error())
		return nil, err
	}
	defer l.Close()
	searchReq := ldap.NewSearchRequest(PeopleDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, "(cn=*)", nil, nil)

	result, err := l.Search(searchReq)

	if err != nil {
		log.Println("QueryRBACRoles: " + err.Error())
		return nil, err
	}
	return result.Entries, err
}

func QueryUser(user string) ([]*ldap.Entry, error) {
	l, err := LDAPConn()
	if err != nil {
		log.Println("QueryUser: " + err.Error())
		return nil, err
	}
	defer l.Close()
	filter := fmt.Sprintf("(cn=%s)", user)
	searchReq := ldap.NewSearchRequest(PeopleDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, filter, nil, nil)

	result, err := l.Search(searchReq)

	if err != nil {
		log.Println("QueryUser: " + err.Error())
		return nil, err
	}
	return result.Entries, err
}

func QueryUsers() ([]*ldap.Entry, error) {
	//This is the same as QueryRBACRoles, same DN
	l, err := LDAPConn()
	if err != nil {
		log.Println("QueryUser: " + err.Error())
		return nil, err
	}
	defer l.Close()
	searchReq := ldap.NewSearchRequest(PeopleDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, "(cn=*)", nil, nil)

	result, err := l.Search(searchReq)

	if err != nil {
		log.Println("QueryUser: " + err.Error())
		return nil, err
	}
	return result.Entries, err
}

func AddAdminUser() error {
	l, err := LDAPConn()

	if err != nil {
		fmt.Println("AddAdminUser: " + err.Error())
		return err
	}
	defer l.Close()
	adminAccountName := "admin"
	dn := fmt.Sprintf("cn=%s,%s", adminAccountName, PeopleDN)
	addReq := ldap.NewAddRequest(dn, nil)
	addAttributeToAddRequest(addReq, "objectClass", []string{"inetOrgPerson", "top", "shadowAccount"})
	addAttributeToAddRequest(addReq, "userPassword", []string{adminAccountName})
	addAttributeToAddRequest(addReq, "sn", []string{adminAccountName})
	addAttributeToAddRequest(addReq, "uid", []string{adminAccountName})
	if err := l.Add(addReq); err != nil {
		log.Println("AddAdminUser: ", err.Error())
		return err
	}
	return nil
}

func AddRBACRole(role string) error {
	l, err := LDAPConn()

	if err != nil {
		fmt.Println("addRBACRole: " + err.Error())
		return err
	}
	defer l.Close()
	dn := fmt.Sprintf("cn=%s,%s", role, PeopleDN)
	addReq := ldap.NewAddRequest(dn, nil)
	addAttributeToAddRequest(addReq, "objectClass", []string{"groupOfUniqueNames", "top"})
	addAttributeToAddRequest(addReq, "cn", []string{role})
	addAttributeToAddRequest(addReq, "uniqueMember", []string{"cn=admin,ou=People,dc=local"})
	if err := l.Add(addReq); err != nil {
		log.Println("addRBACRole: ", addReq, err)
	}

	return nil
}

func DeleteRBACRole(role string) error {
	dn := fmt.Sprintf("cn=%s,%s", role, PeopleDN)
	return delete(dn)
}

func AddUserToRBACRole(name string, role RBACRole) error {
	l, err := LDAPConn()
	if err != nil {
		fmt.Println("AddUserToRBACRole: " + err.Error())
		return err
	}
	defer l.Close()
	var userDNs []string
	dn := fmt.Sprintf("cn=%s,%s", name, PeopleDN)
	userDNs = append(userDNs, dn)
	//userDNs = ensureAdminDN(userDNs)
	dnRole := fmt.Sprintf("cn=%s,%s", role.Name, PeopleDN)
	mod := ldap.NewModifyRequest(dnRole, nil)
	mod.Add("uniqueMember", userDNs)
	if err := l.Modify(mod); err != nil {
		log.Println("AddUserToRBACRole: " + err.Error())
		return err
	}
	return nil
}

func RemoveUserFromRBACRole(name string, role RBACRole) error {
	l, err := LDAPConn()
	if err != nil {
		fmt.Println("RemoveUserFromRBACRole: " + err.Error())
		return err
	}
	defer l.Close()
	var userDNs []string
	for _, r := range role.Users {
		if r.Username == name {
			continue
		}
		dn := fmt.Sprintf("cn=%s,%s", r.Username, PeopleDN)
		userDNs = append(userDNs, dn)
	}
	userDNs = ensureAdminDN(userDNs)
	dn := fmt.Sprintf("cn=%s,%s", role.Name, PeopleDN)
	mod := ldap.NewModifyRequest(dn, nil)
	mod.Replace("uniqueMember", userDNs)
	if err := l.Modify(mod); err != nil {
		log.Println("RemoveUserFromRBACRole: " + err.Error())
		return err
	}
	return nil
}

func ensureAdminDN(userDNs []string) []string {
	var found bool
	adminDN := fmt.Sprintf("cn=%s,%s", "admin", PeopleDN)
	for _, u := range userDNs {
		if u == adminDN {
			found = true
		}
	}
	if !found {
		userDNs = append(userDNs, adminDN)
	}
	return userDNs
}

func AddRBACClaim(name string, value int) error {
	l, err := LDAPConn()
	if err != nil {
		fmt.Println("AddRBACClaim: " + err.Error())
		return err
	}
	defer l.Close()
	dn := fmt.Sprintf("cn=claim-%s-%d,%s", name, value, ClaimDN)
	addReq := ldap.NewAddRequest(dn, nil)
	addAttributeToAddRequest(addReq, "objectClass", []string{"organizationalRole", "x-cluster-claim", "top"})
	addAttributeToAddRequest(addReq, "x-cluster-claim-name", []string{name})
	addAttributeToAddRequest(addReq, "x-cluster-claim-value", []string{fmt.Sprint(value)})
	uniquePtr := fmt.Sprintf("cluster-claim-ptr-%s-%d", name, value)
	addAttributeToAddRequest(addReq, "x-cluster-claim-unique-ptr", []string{uniquePtr})
	addAttributeToAddRequest(addReq, "cn", []string{"cluster"})
	if err := l.Add(addReq); err != nil {
		log.Println("AddRBACClaim: " + err.Error())
		return err
	}
	return nil
}

func DeleteRBACClaim(name string, value int) error {
	dn := fmt.Sprintf("cn=claim-%s-%d,%s", name, value, ClaimDN)
	return delete(dn)
}

func AddRBACClaimToRole(role RBACRole, claimDN string) error {
	//cn=TestRole,ou=People,dc=local
	l, err := LDAPConn()
	if err != nil {
		fmt.Println("AddRBACClaimToRole: " + err.Error())
		return err
	}
	defer l.Close()
	dn := fmt.Sprintf("cn=cluster-claims,cn=%s,%s", role.Name, PeopleDN)
	e, err := QueryRBACRoleClaim(role.Name)
	if err != nil {
		log.Println("AddClaimToRBACRole: " + err.Error())
		return err
	}
	claims := RoleEntriesToRBACClaims(e)
	if len(claims) == 0 {
		addReq := ldap.NewAddRequest(dn, nil)
		addAttributeToAddRequest(addReq, "objectClass", []string{"groupOfUniqueNames", "top"})
		addAttributeToAddRequest(addReq, "uniqueMember", []string{claimDN})
		addAttributeToAddRequest(addReq, "cn", []string{"cluster-claims"})
		//dnRole := fmt.Sprintf("cn=%s,%s", role.Name, PeopleDN)
		if err := l.Add(addReq); err != nil {
			log.Println("AddRBACClaimToRole: " + err.Error())
			return err
		}
	} else {
		mod := ldap.NewModifyRequest(dn, nil)
		mod.Add("uniqueMember", []string{claimDN})
		if err := l.Modify(mod); err != nil {
			log.Println("AddClaimToRBACRole: " + err.Error())
			return err
		}
		/*if err := l.Add(addReq); err != nil {
			log.Println("AddRBACClaimToRole: " + err.Error())
			return err
		}*/
	}
	return nil
}

func DeleteRBACClaimFromRole(roleName string, claimDN string) error {
	l, err := LDAPConn()
	if err != nil {
		fmt.Println("DeleteRBACClaimFromRole: " + err.Error())
		return err
	}
	defer l.Close()
	dn := fmt.Sprintf("cn=cluster-claims,cn=%s,%s", roleName, PeopleDN)
	e, err := QueryRBACRoleClaim(roleName)
	if err != nil {
		log.Println("DeleteRBACClaimFromRole: " + err.Error())
		return err
	}
	claims := RoleEntriesToRBACClaims(e)
	var updatedClaimDNs []string
	for _, c := range claims {
		if c.DN != claimDN {
			updatedClaimDNs = append(updatedClaimDNs, c.DN)
		}
	}

	if len(updatedClaimDNs) == 0 {
		deleteDN(dn)
	} else {
		mod := ldap.NewModifyRequest(dn, nil)
		mod.Replace("uniqueMember", updatedClaimDNs)
		if err := l.Modify(mod); err != nil {
			log.Println("DeleteRBACClaimFromRole: " + err.Error())
			return err
		}
	}

	return nil
}

func addAttributeToAddRequest(req *ldap.AddRequest, typeString string, values []string) {
	if req.Attributes == nil {
		req.Attributes = []ldap.Attribute{}
	}
	attr := ldap.Attribute{
		Type: typeString,
		Vals: values,
	}
	req.Attributes = append(req.Attributes, attr)
}

func DeleteClaim(name string, value int) error {
	dn := fmt.Sprintf("cn=claim-%s-%d,%s", name, value, ClaimDN)
	return delete(dn)
}

func delete(dn string) error {
	l, err := LDAPConn()
	if err != nil {
		fmt.Println("delete: " + err.Error())
		return err
	}
	defer l.Close()

	delReq := ldap.NewDelRequest(dn, nil)
	err = l.Del(delReq)
	if err != nil {
		log.Println("delete: " + err.Error())
		return err
	}
	return nil
}
