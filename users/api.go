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
	r.POST("/validateToken", validateToken)
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
	if Validate(login) {
		cnx.JSON(200, "")
		return
	}

	// case: Invalid credentials
	cnx.JSON(401, "Invalid login credentials")
}

func validateToken(cnx *gin.Context) {

	jwtToken, err := extractBearerToken(cnx.GetHeader("Authorization"))
	if err != nil {
		cnx.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: err.Error(),
		})
		return
	}

	token, err := ValidateToken(jwtToken)
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
	cnx.JSON(200, "")
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
