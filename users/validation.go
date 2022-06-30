package users

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func validateToken(cnx *gin.Context) error {
	cookie, err := cnx.Request.Cookie(AccessTokenKey)
	if err != nil {
		return err
	}

	token, err := ValidateJWTToken(cookie.Value)
	if err != nil {
		if err.Error() == "Token is expired" {
			claims, OK := token.Claims.(jwt.MapClaims)
			if !OK {
				return err
			}
			username := claims["username"].(string)
			err := refreshExpiredToken(cnx, username)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	_, OK := token.Claims.(jwt.MapClaims)
	if !OK {
		return err
	}
	return nil
}

func refreshExpiredToken(cnx *gin.Context, username string) error {
	cookie, err := cnx.Request.Cookie(RefreshTokenKey)
	if err != nil {
		return err
	}
	baseDN := "ou=People,DC=local"
	ldapRefreshToken, err := queryLDAPUserAttribute(baseDN, username, "description")
	if err != nil {
		return err
	}
	if ldapRefreshToken == cookie.Value {
		userInfo := UserInfo{Username: username}
		if cookie.Value == ldapRefreshToken {
			jwtString, rtString, err := RefreshJWTToken(ldapRefreshToken, userInfo)
			expirationTime := time.Now().Add(AccessTokenExpirationWindow)
			refreshExpirationTime := time.Now().Add(RefreshTokenExpirationWindow)
			cnx.SetCookie(AccessTokenKey, jwtString, int(expirationTime.Unix()), "/", "localhost", false, true)
			cnx.SetCookie(RefreshTokenKey, rtString, int(refreshExpirationTime.Unix()), "/", "localhost", false, true)
			return err
		}
	}
	return errors.New("invalid refresh token")
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
