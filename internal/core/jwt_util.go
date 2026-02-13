package core

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwtToken(userId, scr string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	})

	secret := []byte(scr)
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to create token")
	}

	return tokenStr, nil
}

func ValidateJwtToken(tokenStr, scr string) (*jwt.Token, error) {
	secret := []byte(scr)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, InternalServerError("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		return nil, BadRequest("invalid token: " + err.Error())
	}

	if !token.Valid {
		return nil, BadRequest("invalid token")
	}

	return token, nil
}
