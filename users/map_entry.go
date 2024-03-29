package users

import (
	"log"
	"strconv"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

func EntriesToRBACClaims(entries []*ldap.Entry) (claims []RBACClaim) {
	for _, e := range entries {
		claimName := e.GetAttributeValue("x-cluster-claim-name")
		if claimName != "" {
			var rbacClaim RBACClaim
			rbacClaim.Name = claimName
			valString := e.GetAttributeValue("x-cluster-claim-value")
			val, err := strconv.Atoi(valString)
			if err != nil {
				log.Println("EntriesToRBACClaims: " + err.Error())
				continue
			}
			rbacClaim.Value = val
			rbacClaim.DN = e.DN
			claims = append(claims, rbacClaim)
		}
	}
	return
}

func RoleEntriesToRBACClaims(entries []*ldap.Entry) (claims []RBACClaim) {
	for _, e := range entries {
		cn := e.GetAttributeValue("cn")
		if cn != "cluster-claims" {
			continue
		}
		res := e.GetAttributeValues("uniqueMember")
		for _, r := range res {
			var rbacClaim RBACClaim
			rbacClaim.Name, rbacClaim.Value = parseClusterClaimCN(r)
			rbacClaim.DN = r
			claims = append(claims, rbacClaim)
		}
	}
	return claims
}

func EntriesToRoles(entries []*ldap.Entry) (roles []RBACRole) {
	for _, e := range entries {
		class := e.GetAttributeValue("objectClass")
		if class == "groupOfUniqueNames" {
			var role RBACRole
			role.Name = e.GetAttributeValue("cn")
			if role.Name == "cluster-claims" {
				continue
			}
			userStrings := e.GetAttributeValues("uniqueMember")
			for _, u := range userStrings {
				var user UserInfo
				user.Username = parseUserCN(u)
				role.Users = append(role.Users, user)
			}
			roles = append(roles, role)
		}
	}
	return
}

func RoleEntriesToUsers(entries []*ldap.Entry) (users []UserInfo) {
	for _, e := range entries {
		class := e.GetAttributeValue("objectClass")
		if class == "groupOfUniqueNames" && !strings.HasPrefix(e.DN, "cn=cluster-claims") {
			vals := e.GetAttributeValues("uniqueMember")
			for _, v := range vals {
				var user UserInfo
				user.Username = parseUserCN(v)
				users = append(users, user)
			}
		}
	}
	return
}

func EntriesToUsers(entries []*ldap.Entry) (users []UserInfo) {
	for _, e := range entries {
		uid := e.GetAttributeValue("uid")
		if uid != "" {
			var user UserInfo
			user.Username = parseUserCN(e.DN)
			users = append(users, user)
		}
	}
	return
}

func parseUserCN(userCN string) (name string) {
	s := strings.Split(userCN, ",")
	if len(s) == 0 {
		log.Printf("parseUserCN: ERROR - invalid user CN %s", userCN)
		return
	}
	a := strings.Split(s[0], "=")
	if len(a) == 0 {
		log.Printf("parseUserCN: ERROR - invalid user CN %s", userCN)
		return
	}
	name = a[1]
	return
}

func parseClusterClaimCN(clusterClaimCN string) (name string, value int) {
	//"cn=claim-test2-100,cn=cluster,dc=local"
	s := strings.Split(clusterClaimCN, ",")
	if len(s) == 0 {
		log.Printf("parseClusterClaimCN: ERROR - no cluster claims found in %s", clusterClaimCN)
		return
	}
	a := strings.Split(s[0], "=")
	if len(a) == 0 {
		log.Printf("parseClusterClaimCN: ERROR - invalid CN in %s", clusterClaimCN)
		return
	}
	c := strings.Split(a[1], "-")
	if len(c) == 0 {
		log.Printf("parseClusterClaimCN: ERROR - invalid claim in %s", clusterClaimCN)
		return
	}
	name = c[1]
	value, err := strconv.Atoi(c[2])
	if err != nil {
		log.Printf("parseClusterClaimCN: ERROR - invalid claim value in %s", clusterClaimCN)
		return
	}
	return
}
