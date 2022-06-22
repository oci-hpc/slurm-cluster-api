package users

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

//jwt service
type JWTService interface {
	GenerateToken(userInfo UserInfo) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

func getSecretKey() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// default secret; TODO: remove and cleanup env vars
		secret = "WRE8l0A6FUQhZ8FKPzp9Vx0Jg0ANpCt1"
	}
	return secret
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return getSecretKey(), nil
	})

	return jwtToken, err
}

func GenerateToken(userInfo UserInfo) (string, error) {
	return "fake", nil
}
