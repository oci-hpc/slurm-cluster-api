package users

import (
	"errors"
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
	r.GET("/validateToken", validateTokenEndpoint)
	r.POST("/refreshToken", refreshToken)
}

func login(cnx *gin.Context) {
	// pluck user and password
	var login LoginInfo
	if err := cnx.ShouldBind(&login); err != nil {
		cnx.JSON(400, errors.New("missing login information"))
		return
	}

	// check vs ldap
	if userInfo, ok := validateLDAPLogin(login); ok {
		jwtToken, refreshToken, err := GenerateJWTToken(userInfo)
		if err != nil {
			cnx.JSON(500, "could not generate a token")
		}
		expirationTime := time.Now().Add(5 * time.Minute)
		refreshExpirationTime := time.Now().Add(7 * 24 * time.Hour)
		cnx.SetCookie("access_token", jwtToken, int(expirationTime.Unix()), "", "localhost", false, true)
		cnx.SetCookie("refreshToken", refreshToken, int(refreshExpirationTime.Unix()), "/", "localhost", false, true)
		cnx.JSON(200, "")
		return
	}

	// case: Invalid credentials
	cnx.JSON(401, "Invalid login credentials")
}

func validateTokenEndpoint(cnx *gin.Context) {
	validateToken(cnx)
}

func refreshToken(cnx *gin.Context) {
	cnx.JSON(200, "")
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
