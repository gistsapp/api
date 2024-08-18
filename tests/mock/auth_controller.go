package mock

import (
	"github.com/gistapp/api/auth"
	"github.com/gofiber/fiber/v2"
)

type MockAuthController struct{
	AuthService auth.IAuthService
}

func (a *MockAuthController) Callback() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (a *MockAuthController) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (a *MockAuthController) LocalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		e := new(auth.AuthLocalValidator)
		if err := c.BodyParser(e); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with field email as text")
		}
		token, err := a.AuthService.LocalAuth(e.Email)

		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		return c.JSON(fiber.Map{"token": token.Value.String})
	}
}

func (a *MockAuthController) VerifyAuthToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		e := new(auth.AuthLocalVerificationValidator)

		if err := c.BodyParser(e); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with fields token and email as text")
		}

		token := e.Token
		email := e.Email

		jwt_token, err := a.AuthService.VerifyLocalAuthToken(token, email)

		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		token_cookie := new(fiber.Cookie)
		token_cookie.Name = "gists.access_token"
		token_cookie.HTTPOnly = true
		token_cookie.Value = jwt_token
		c.Cookie(token_cookie)
		return c.Status(200).JSON(fiber.Map{"message": "You are now logged in"})
	}
}