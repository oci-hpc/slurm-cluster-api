package users

import (
	"log"
	"strconv"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

func EntriesToRBACClaims(entries []*ldap.Entry) (claims []RBACClaim) {
	for _, e := range entries {
		cn := e.GetAttributeValue("cn")
		if cn != "cluster-claims" {
			return claims
		}
		res := e.GetAttributeValues("uniqueMember")
		for _, r := range res {
			var rbacClaim RBACClaim
			rbacClaim.Name, rbacClaim.Value = parseClusterClaimCN(r)
			claims = append(claims, rbacClaim)
		}
	}
	return claims
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
