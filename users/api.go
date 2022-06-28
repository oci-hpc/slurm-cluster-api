package users

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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
	r.POST("/validateToken", validateToken)
	r.POST("/refreshToken", TokenAuthMiddleware(), refreshToken)
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

		tokens := map[string]string{
			"access_token":  jwtToken,
			"refresh_token": refreshToken,
		}
		cnx.JSON(http.StatusCreated, tokens)

		return
	}

	// case: Invalid credentials
	cnx.JSON(http.StatusUnauthorized, "Invalid login credentials")

}

func logout(cnx *gin.Context) {

	// revoke token from store
	// jwtToken, err := extractBearerToken(cnx.GetHeader("Authorization"))
	// if err != nil {
	// 	cnx.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
	// 		Message: err.Error(),
	// 	})
	// 	return
	// }
	// revokeRefreshToken(refreshToken)

	// case: Invalid credentials
	cnx.JSON(http.StatusUnauthorized, "Invalid login credentials")
}

func validateToken(cnx *gin.Context) {

	jwtToken, err := extractBearerToken(cnx.GetHeader("Authorization"))
	if err != nil {
		cnx.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: err.Error(),
		})
		return
	}

	token, err := ValidateJWTToken(jwtToken)
	if err != nil {
		cnx.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: "bad jwt token",
		})
		return
	}

	_, OK := token.Claims.(jwt.MapClaims)
	if !OK {
		cnx.AbortWithStatusJSON(http.StatusInternalServerError, UnsignedResponse{
			Message: "unable to parse claims",
		})
		return
	}
	cnx.Next()
}

func refreshToken(cnx *gin.Context) {
	mapToken := map[string]string{}
	if err := cnx.ShouldBindJSON(&mapToken); err != nil {
		cnx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	oldRefreshToken := mapToken["refresh_token"]
	userInfo := UserInfo{}
	tokenString, refreshTokenString, err := RefreshJWTToken(oldRefreshToken, userInfo)
	if err != nil {
		cnx.JSON(http.StatusUnauthorized, "Invalid login credentials")
		return
	}

	tokens := map[string]string{
		"access_token":  tokenString,
		"refresh_token": refreshTokenString,
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

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}
