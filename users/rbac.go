package users

import (
	"errors"
	"fmt"
	"log"

	ldap "github.com/go-ldap/ldap/v3"
)

func CheckRoleExists(role string) bool {
	res, err := queryRBACRole(role)
	if err != nil {
		return false
	}
	if len(res) == 0 {
		return false
	}
	return true
}

func queryRBACRole(role string) ([]*ldap.EntryAttribute, error) {
	l, err := LDAPConn()
	if err != nil {
		log.Println("queryRBACRole: " + err.Error())
		return nil, err
	}
	defer l.Close()
	// Filters must start and finish with ()
	filter := fmt.Sprintf("(cn=%s)", ldap.EscapeFilter(role))

	searchReq := ldap.NewSearchRequest(BaseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, filter, nil, nil)

	result, err := l.Search(searchReq)
	if err != nil {
		log.Println("queryRBACRole: " + err.Error())
		return nil, err
	}
	if len(result.Entries) == 0 {
		err = errors.New("Could not find role")
		log.Println("queryRBACRole: " + err.Error())
		return nil, err
	}
	return result.Entries[0].Attributes, err
}

func AddAdminUser() error {
	l, err := LDAPConn()

	if err != nil {
		fmt.Println("AddAdminUser: " + err.Error())
		return err
	}
	defer l.Close()
	adminAccountName := "admin"
	dn := fmt.Sprintf("cn=%s,ou=People,dc=local", adminAccountName)
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
	dn := fmt.Sprintf("cn=%s,ou=People,dc=local", role)
	addReq := ldap.NewAddRequest(dn, nil)
	addAttributeToAddRequest(addReq, "objectClass", []string{"groupOfUniqueNames", "top"})
	addAttributeToAddRequest(addReq, "cn", []string{role})
	addAttributeToAddRequest(addReq, "uniqueMember", []string{"cn=admin,ou=People,dc=local"})
	if err := l.Add(addReq); err != nil {
		log.Println("addRBACRole: ", addReq, err)
	}

	return nil
}

func AddRBACClaim(name string, value int32) error {
	l, err := LDAPConn()
	if err != nil {
		fmt.Println("AddRBACClaim: " + err.Error())
		return err
	}
	defer l.Close()
	dn := fmt.Sprintf("cn=claim-%s-%d,cn=cluster,dc=local", name, value)
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

func AddRBACClaimToRole(role string, claimDN string) error {
	//cn=TestRole,ou=People,dc=local
	l, err := LDAPConn()
	if err != nil {
		fmt.Println("AddRBACClaimToRole: " + err.Error())
		return err
	}
	defer l.Close()
	dn := fmt.Sprintf("cn=cluster-claims,cn=%s,%s", role, BaseDN)
	addReq := ldap.NewAddRequest(dn, nil)
	addAttributeToAddRequest(addReq, "objectClass", []string{"groupOfUniqueNames", "top"})
	addAttributeToAddRequest(addReq, "uniqueMember", []string{claimDN})
	addAttributeToAddRequest(addReq, "cn", []string{"cluster-claims"})
	if err := l.Add(addReq); err != nil {
		log.Println("AddRBACClaimToRole: " + err.Error())
		return err
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
