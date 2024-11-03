package mock

import (
	"github.com/gistapp/api/user"
	"github.com/gistapp/api/utils"
	"github.com/gofiber/fiber/v2"
)

type MockAuthController struct {
	AuthService user.IAuthService
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
		e := new(user.AuthLocalValidator)
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

func (a *MockAuthController) Renew() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user_id := c.Locals("pub").(string)

		tokens, err := a.AuthService.Renew(user_id)

		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		c.Cookie(utils.Cookie("gists.access_token", tokens.AccessToken))   //set access token
		c.Cookie(utils.Cookie("gists.refresh_token", tokens.RefreshToken)) //set refresh token

		return c.Status(200).JSON(fiber.Map{"message": "Welcome back"})
	}
}

func (a *MockAuthController) VerifyAuthToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		e := new(user.AuthLocalVerificationValidator)

		if err := c.BodyParser(e); err != nil {
			return c.Status(400).SendString("Request must be valid JSON with fields token and email as text")
		}

		token := e.Token
		email := e.Email

		jwt_token, err := a.AuthService.VerifyLocalAuthToken(token, email)

		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		c.Cookie(utils.Cookie("gists.access_token", jwt_token.AccessToken))   //set access token
		c.Cookie(utils.Cookie("gists.refresh_token", jwt_token.RefreshToken)) //set refresh token
		return c.Status(200).JSON(fiber.Map{"message": "You are now logged in"})
	}
}

func (a *MockAuthController) Logout() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.ClearCookie("gists.access_token")
		c.ClearCookie("gists.refresh_token")
		return c.Status(200).JSON(fiber.Map{"message": "You are now logged out"})
	}
}
