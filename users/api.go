package users

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UnsignedResponse struct {
	Message interface{} `json:"message"`
}

type SignedResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

func InitializeUsersEndpoint(r *gin.Engine) {
	r.POST("/login", login)
	r.POST("/logout", TokenAuthMiddleware(), logout)
	r.GET("/validateToken", validateTokenEndpoint)
	r.GET("/logout", logout)
	r.POST("/refreshToken", TokenAuthMiddleware(), refreshToken)
	r.GET("/claims", getClaims)
	r.POST("/claims", addClaim)
}

func login(cnx *gin.Context) {
	// pluck user and password
	var login LoginInfo
	if err := cnx.ShouldBind(&login); err != nil {
		cnx.JSON(http.StatusBadRequest, errors.New("missing login information"))
		return
	}

	// check vs ldap
	if userInfo, ok := validateLDAPLogin(login); ok {
		jwtToken, refreshToken, err := GenerateJWTToken(userInfo)
		if err != nil {
			cnx.JSON(http.StatusInternalServerError, "could not generate a token")
		}
		expirationTime := time.Now().Add(AccessTokenExpirationWindow)
		refreshExpirationTime := time.Now().Add(RefreshTokenExpirationWindow)
		cnx.SetCookie(AccessTokenKey, jwtToken, int(expirationTime.Unix()), "/", "localhost", false, true)
		cnx.SetCookie(RefreshTokenKey, refreshToken, int(refreshExpirationTime.Unix()), "/", "localhost", false, true)
		cnx.JSON(200, userInfo)
		return
	}

	// case: Invalid credentials
	cnx.JSON(http.StatusUnauthorized, "Invalid login credentials")

}

func logout(cnx *gin.Context) {
	cnx.SetCookie(AccessTokenKey, "", 0, "/", "localhost", false, true)
	cnx.SetCookie(RefreshTokenKey, "", 0, "/", "localhost", false, true)
	cnx.Status(200)
}

func validateTokenEndpoint(cnx *gin.Context) {
	err := validateToken(cnx)
	if err != nil {
		cnx.AbortWithStatusJSON(http.StatusUnauthorized, UnsignedResponse{
			Message: "Unauthorized",
		})
	}
	cnx.Next()
}

func refreshToken(cnx *gin.Context) {
	mapToken := map[string]string{}
	if err := cnx.ShouldBindJSON(&mapToken); err != nil {
		cnx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	oldRefreshToken := mapToken[RefreshTokenKey]
	userInfo := UserInfo{}
	tokenString, refreshTokenString, err := RefreshJWTToken(oldRefreshToken, userInfo)
	if err != nil {
		cnx.JSON(http.StatusUnauthorized, "Invalid login credentials")
		return
	}

	tokens := map[string]string{
		AccessTokenKey:  tokenString,
		RefreshTokenKey: refreshTokenString,
	}
	cnx.JSON(http.StatusCreated, tokens)
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(cnx *gin.Context) {
		token, err := extractBearerToken(cnx.GetHeader("Authorization"))
		if err != nil {
			cnx.JSON(http.StatusUnauthorized, err.Error())
			cnx.Abort()
			return
		}

		_, err = ValidateJWTToken(token)
		if err != nil {
			cnx.JSON(http.StatusUnauthorized, err.Error())
			cnx.Abort()
			return
		}
		cnx.Next()
	}
}

func getClaims(cnx *gin.Context) {
	query := cnx.Request.URL.Query()
	var res []RBACClaim
	if val, ok := query["role"]; ok {
		entries, err := QueryRBACRoleClaim(val[0])
		if err != nil {
			cnx.JSON(http.StatusInternalServerError, "Server internal error")
		}
		res = EntriesToRBACClaims(entries)
		cnx.JSON(http.StatusOK, res)
	} else {
		entries, err := QueryAllRBACClaims()
		if err != nil {
			cnx.JSON(http.StatusInternalServerError, "Server internal error")
		}
		res = EntriesToRBACClaims(entries)
		cnx.JSON(http.StatusOK, res)
	}
}

func addClaim(cnx *gin.Context) {
	var claim RBACClaim
	if err := cnx.ShouldBindJSON(&claim); err != nil {
		cnx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	if claim.Name == "" {
		cnx.JSON(http.StatusBadRequest, "Name is invalid")
		return
	}
	query := cnx.Request.URL.Query()
	if val, ok := query["role"]; ok {
		e, err := QueryRBACClaim(claim.Name, claim.Value)
		if err != nil {
			log.Printf("addClaim - %s", err.Error())
			cnx.JSON(http.StatusBadRequest, "Invalid claim")
			return
		}
		if len(e) == 0 {
			log.Printf("addClaim - no claims found in query that match LDAP claims")
			cnx.JSON(http.StatusBadRequest, "Invalid claim")
			return
		}
		err = AddRBACClaimToRole(val[0], e[0].DN)
		if err != nil {
			log.Printf("addClaim - %s", err.Error())
			cnx.JSON(http.StatusBadRequest, "Unable to add claim to role")
			return
		}
	} else {
		err := AddRBACClaim(claim.Name, claim.Value)
		if err != nil {
			cnx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		cnx.JSON(http.StatusOK, "")
	}
}
