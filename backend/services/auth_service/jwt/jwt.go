package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secret []byte

func Init(key string) {
	secret = []byte(key)
}

func GenerateToken(userID string) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ValidateToken(tkn string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tkn, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
