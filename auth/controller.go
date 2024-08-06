package auth

import (
	"github.com/gofiber/fiber/v2"
)

type AuthControllerImpl struct{}

type AuthLocalValidator struct {
	Email string `json:"email"`
}

type AuthLocalVerificationValidator struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

func (a *AuthControllerImpl) Callback() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return AuthService.Callback(c)
	}
}

func (a *AuthControllerImpl) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return AuthService.Authenticate(c)
	}
}

func (a *AuthControllerImpl) LocalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		e := new(AuthLocalValidator)
		if err := c.BodyParser(e); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with field email as text")
		}

		if err := AuthService.LocalAuth(e.Email); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		return c.JSON(fiber.Map{"message": "Check your email for the token"})
	}
}

func (a *AuthControllerImpl) VerifyAuthToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		e := new(AuthLocalVerificationValidator)

		if err := c.BodyParser(e); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with fields token and email as text")
		}

		token := e.Token
		email := e.Email

		jwt_token, err := AuthService.VerifyLocalAuthToken(token, email)

		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		return c.JSON(fiber.Map{"token": jwt_token})
	}
}

var AuthController AuthControllerImpl = AuthControllerImpl{}
