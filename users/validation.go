package users

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func validateToken(cnx *gin.Context) bool {
	cookie, err := cnx.Request.Cookie("access_token")
	if err != nil {
		cnx.AbortWithStatusJSON(http.StatusUnauthorized, UnsignedResponse{
			Message: "Unauthorized",
		})
		return false
	}

	token, err := ValidateJWTToken(cookie.Value)
	if err != nil {
		cnx.AbortWithStatusJSON(http.StatusUnauthorized, UnsignedResponse{
			Message: "Unauthorized",
		})
		return false
	}

	_, OK := token.Claims.(jwt.MapClaims)
	if !OK {
		cnx.AbortWithStatusJSON(http.StatusUnauthorized, UnsignedResponse{
			Message: "Unauthorized",
		})
		return false
	}
	cnx.Next()
	return true
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
