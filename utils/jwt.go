package utils

import (
	"errors"
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

func VerifyJWT(raw_token string) (map[string]any, error) {
	token, err := jwt.Parse(raw_token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return map[string]any{
				"error": "Unauthorized",
			}, errors.New("unauthorized")
		}
		return secretKey, nil
	})

	if err != nil {
		return map[string]any{}, err
	}

	var toReturn map[string]any = make(map[string]any)

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		for key, val := range claims {
			toReturn[key] = val
		}
		return toReturn, nil
	}
	return map[string]any{}, ErrInvalidToken
}

var ErrInvalidToken = errors.New("jwt token is invalid for some reason")
