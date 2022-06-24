package users

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func getSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// default secret; TODO: remove and cleanup env vars
		secret = "xIKYselVMMc5XS5ATExrD30OuxKTt8F4eMZba62TGVk="
	}
	return []byte(secret)
}

func getRefreshSecretKey() []byte {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		// default secret; TODO: remove and cleanup env vars
		secret = "xIKYselVMMc5XS5ATExrD30OuxKTt8F4eMZba62TGVk="
	}
	return []byte(secret)
}

func ValidateJWTToken(tokenString string) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getSecretKey(), nil
	})

	return jwtToken, err
}

func GenerateJWTToken(userInfo UserInfo) (tokenString string, refreshTokenString string, err error) {
	tokenString, err = generateAccessToken(userInfo)
	if err != nil {
		return "", "", err
	}
	refreshTokenString, err = generateRefreshToken(userInfo)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, err
}

func RefreshJWTToken(refreshToken string, userInfo UserInfo) (tokenString string, refreshTokenString string, err error) {

	//validate refresh token
	_, err = jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getRefreshSecretKey(), nil
	})

	if err != nil {
		return "", "", err
	}

	// if refresh token valid, generate new tokens
	tokenString, err = generateAccessToken(userInfo)
	if err != nil {
		return "", "", err
	}
	refreshTokenString, err = generateRefreshToken(userInfo)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, err
}

func generateAccessToken(userInfo UserInfo) (refreshTokenString string, err error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &JWTClaim{
		Email:    userInfo.Email,
		Username: userInfo.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshTokenString, err = token.SignedString(getSecretKey())
	return refreshTokenString, err
}

func generateRefreshToken(userInfo UserInfo) (refreshTokenString string, err error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &JWTClaim{
		Email:    userInfo.Email,
		Username: userInfo.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshTokenString, err = token.SignedString(getRefreshSecretKey())
	return refreshTokenString, err
}
