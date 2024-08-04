package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(Get("APP_KEY"))

func CreateToken(email string, user_id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"pub":   user_id,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
