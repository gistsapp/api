package utils

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(Get("APP_KEY"))

type AccessToken struct {
	Email string `jwt:"email"`
	Pub   string `jwt:"pub"`
}

type RefreshToken struct {
	Pub string `jwt:"pub"`
}

func CreateRefreshToken(user_id string) (string, error) {
	payload := RefreshToken{
		Pub: user_id,
	}
	return CreateToken(payload)
}

func CreateAccessToken(email string, user_id string) (string, error) {
	payload := AccessToken{
		Email: email,
		Pub:   user_id,
	}
	return CreateToken(payload)
}

func CreateToken(payload interface{}) (string, error) {

	t_payload := reflect.TypeOf(payload)
	v_payload := reflect.ValueOf(payload)

	claims := make(jwt.MapClaims)

	for i := 0; i < t_payload.NumField(); i++ {
		field := t_payload.Field(i)

		tag_value := field.Tag.Get("jwt")
		if tag_value == "" {
			continue
		}

		log.Info(v_payload.Field(i))

		claims[tag_value] = fmt.Sprintf("%s", v_payload.Field(i))
	}

	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	log.Info(claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		claims,
	)

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

	var to_return map[string]any = make(map[string]any)

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		for key, val := range claims {
			to_return[key] = val
		}
		return to_return, nil
	}
	return map[string]any{}, ErrInvalidToken
}

var ErrInvalidToken = errors.New("jwt token is invalid for some reason")
